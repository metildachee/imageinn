import json

def read_json_file_line_by_line(file_path):
    data = []
    with open(file_path, 'r') as file:
        for line in file:
            data.append(json.loads(line))
    return data

def read_category_file(file_path):
    category_map = {}
    with open(file_path, 'r') as file:
        idx = 0
        for line in file:
            idx += 1
            category_name = line.strip()  # Assuming each line is "category_id,category_name"
            category_map[idx] = category_name
    return category_map

# File paths
board_to_category_path = 'board_to_category.json'
board_to_pin_path = 'board_to_pins.json'
pin_to_img_path = 'pin_to_img.json'
category_path = 'categories.txt'

# Reading the files
board_to_category = read_json_file_line_by_line(board_to_category_path)
board_to_pin = read_json_file_line_by_line(board_to_pin_path)
pin_to_img = read_json_file_line_by_line(pin_to_img_path)
category_map = read_category_file(category_path)

print(category_map)

# Map boards to categories
board_category_map = {item["board_id"]: item["cate_id"] for item in board_to_category}

# Map pins to boards
pin_board_map = {}
for item in board_to_pin:
    for pin in item["pins"]:
        pin_board_map[pin] = item["board_id"]

# Create a dictionary to map pin_id to category_id and image URL
pin_details_map = {}
for pin in pin_to_img:
    pin_id = pin["pin_id"]
    img_url = pin["im_url"]
    if pin_id in pin_board_map:
        board_id = pin_board_map[pin_id]
        if board_id in board_category_map:
            cate_id = board_category_map[board_id]
            pin_details_map[pin_id] = {
                'id': pin_id,
                'category_ids': [int(cate_id)],
                'image_url': img_url,
                'caption': category_map[int(cate_id)] 
            }

output_file_path = 'pinterest_es.json'  # Change this to the desired path

with open(output_file_path, 'w') as file:
    for pin_id, details in pin_details_map.items():
        # Creating the action/metadata line
        action = {
            "index": {
                "_index": "images",  # Your index name
                "_id": details['id']  # Use the pin_id as the document ID in Elasticsearch
            }
        }
        # Write the action line
        file.write(json.dumps(action) + '\n')
        
        # Remove the id from the details as it's already used in the index action
        del details['id']

        # Write the document line
        file.write(json.dumps(details) + '\n')