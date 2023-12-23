/**
 * Copyright 2023 Google LLC
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

'use strict'

// Import
const uuid4 = require('uuid4')
const moment = require('moment')

// SCC client
const { SecurityCenterClient } = require('@google-cloud/security-center');
const client = new SecurityCenterClient();

// Environment variables
const sourceId = process.env.SOURCE_ID;

// Variables
const sccConfig = {
    category: 'GCS Bucket - Set IAM policy Alert',
    findingClass: 'OBSERVATION',
    severity: 'MEDIUM'
};

// Exported function
exports.logMonitoring = message => {
    try {
        var event = parseMessage(message);

        // This validate is specific for the Log Monitoring scenario.
        validateEvent(event);

        // Map of extra properties to save on the finding with field name and value
        var extraProperties = {
            description:  event.incident.summary
        };

        createFinding(
            event.incident.started_at,
            event.incident.resource.labels.bucket_name,
            extraProperties
        );
    } catch (error) {
        console.warn(`Skipping executing with message: ${error.message}`);
    }
}

/**
 * Parse the message received on the Cloud Function to a JSON.
 *
 * @param {any} message  Message from Cloud Function
 * @returns {JSON} Json object from the message
 * @exception If some error happens while parsing, it will log the error and finish the execution
 */
function parseMessage(message) {
    // If message data is missing, log a warning and exit.
    if (!(message && message.data)) {
        throw new Error(`Missing required fields (message or message.data)`);
    }

    // Extract the event data from the message
    var event = JSON.parse(Buffer.from(message.data, 'base64').toString());

    return event;
}

/**
 * Validate if the asset is from Organizations and have the iamPolicy and bindings field.
 *
 * @param {any} asset  Asset JSON.
 * @exception If the asset is not valid it will throw the corresponding error.
 */
function validateEvent(event) {
    // If this is not a set-bucket-iam-policy incident, throw an error.
    if (!(event.incident && event.incident.metric && event.incident.metric.type && event.incident.metric.type === "logging.googleapis.com/user/set-bucket-iam-policy")) {
        throw new Error(`Not a set-bucket-iam-policy incident`);
    }
}

/**
 * Create the new SCC finding
 *
 * @param {string} updateTime The time that the asset was changed.
 * @param {string} resourceName The resource where the role was given.
 * @param {Any} extraProperties A key/value map with properties to save on the finding ({fieldName: fieldValue})
 */
async function createFinding(updateTime, resourceName, extraProperties) {
    const [newFinding] = await client.createFinding(
        {
            parent: sourceId,
            findingId: uuid4().replace(/-/g, ''),
            finding: {
                ... {
                    state: 'ACTIVE',
                    resourceName: resourceName,
                    category: sccConfig.category,
                    eventTime: {
                        seconds: updateTime,
                        nanos: 0
                    },
                    findingClass: sccConfig.findingClass,
                    severity: sccConfig.severity
                },
                ...extraProperties
            }
        }
    );

    console.log('New finding created: %j', newFinding);
}
