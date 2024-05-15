import os
import argparse

parser = argparse.ArgumentParser(description="Process some parameters.")
parser.add_argument('--directory', type=str, default='amazon', help='The directory to read bulk and save JSON file.')
args = parser.parse_args()

def split_bulk_json(input_file, output_dir, rows_per_file):
    os.makedirs(output_dir, exist_ok=True)

    with open(input_file, 'r') as file:
        lines = file.readlines()

    total_pairs = len(lines) // 2
    
    num_files = (total_pairs + rows_per_file - 1) // rows_per_file

    for i in range(num_files):
        start = i * rows_per_file
        end = start + rows_per_file
        current_lines = lines[start*2:end*2]

        with open(f'{output_dir}/bulk_part_{i+1}.json', 'w') as output_file:
            output_file.writelines(current_lines)

input_file = args.directory + "/"+ args.directory + "_es.json"
output_dir = args.directory + "/bulk"
rows_per_file = 5000

split_bulk_json(input_file, output_dir, rows_per_file)