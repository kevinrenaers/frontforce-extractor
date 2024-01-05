# Frontforce-extractor

### Configuration file data:

1. **frontforce_url**:

   - Description: URL for accessing the Frontforce application.
   - Value: `https://limburg.frontforce.be/app`

2. **frontforce_username**:

   - Description: Username for accessing the Frontforce application.

3. **frontforce_password**:

   - Description: Password for accessing the Frontforce application.

4. **refresh_interval**:

   - Description: Interval (in seconds) for refreshing data.
   - Value: `5` (Change as per your refresh requirements)

5. **ha_url**:

   - Description: URL for Home Assistant (HA) system.

6. **ha_token**:

   - Description: Authentication token for Home Assistant.

7. **ha_frontforce_status_entity**:

   - Description: Entity information for Frontforce status in Home Assistant.
   - Value:
     ```json
     {
       "id": "input_text.frontforce_status",
       "friendly_name": "frontforce_status"
     }
     ```
   - Update `id` and `friendly_name` as needed for your Home Assistant configuration.

8. **ha_frontforce_intervention_entity**:
   - Description: Entity information for Frontforce intervention in Home Assistant.
   - Value:
     ```json
     {
       "id": "input_text.frontforce_interventie",
       "friendly_name": "frontforce_interventie"
     }
     ```
   - Modify `id` and `friendly_name` based on your Home Assistant setup.

## Prerequisite

Installation of golang

## Usage

1. **_Config file_**

   - Rename the config_example.json to config.json and fill in the needed fields as described above.
   - Place config.json in go bin location
     -> linux: ~/go/bin

2. Compile go program
   ```
   go install
   ```
3. **_Run program_**
   - in go bin directory execute
   ```console
      ./frontforce
   ```
