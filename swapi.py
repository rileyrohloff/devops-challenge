import yaml
import requests as req
import json
import os
import time


SWAPI_URL = "https://swapi.dev/api"
input_file = 'input.yaml'
output_file = 'swapi-output.json'



def get_yaml_input(file_name: str) -> dict:
    try:
        with open(file_name, "r") as f:
            y = yaml.load(f, Loader=yaml.Loader)     
            return y
    except yaml.YAMLError as err:
        print(f"ERROR occured reading yaml file: {err}")

def run_query(swapi_object: str, fields: list) -> dict:
    url = f"{SWAPI_URL}/{swapi_object['type']}/{swapi_object['id']}"
    swapi_dict = {}
    try:
        resp = req.get(url=url)
        resp.raise_for_status()
        data = resp.json()
        for field in fields:
            swapi_dict[field] = data[field]
    except req.exceptions.HTTPError as err:
        print(err)

    return swapi_dict
    


if __name__ == "__main__":
    input_data = get_yaml_input(input_file)
    print("Starting swapi.py script...")
    # keep the pod alive with a while loop :)
    while True:
        ouput_data = []
        f = open(output_file, "w+")
        for item in input_data['input']:
            result = run_query(item, item['infoRequest'])
            ouput_data.append(result)
        json_data = json.dumps(ouput_data, indent=2)
        f.write(json_data)
        print(f'Wrote output to {output_file}')
        f.close()
        time.sleep(30)
        os.remove(output_file)
        print('Removing file....')
    
    