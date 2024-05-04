import os

def split_bulk_json(input_file, output_dir, rows_per_file):
    # Ensure the output directory exists
    os.makedirs(output_dir, exist_ok=True)

    # Open the input file
    with open(input_file, 'r') as file:
        lines = file.readlines()

    # Calculate total pairs (index + data rows are considered a pair)
    total_pairs = len(lines) // 2
    
    # Number of files needed
    num_files = (total_pairs + rows_per_file - 1) // rows_per_file  # ceiling division

    for i in range(num_files):
        # Calculate the slice of lines for the current file
        start = i * rows_per_file
        end = start + rows_per_file
        current_lines = lines[start*2:end*2]  # Multiply by 2 as each pair consists of two lines

        # Write the current slice to a new file
        with open(f'{output_dir}/bulk_part_{i+1}.json', 'w') as output_file:
            output_file.writelines(current_lines)

# Configuration
input_file = 'pinterest_es.json'  # Update this path to your original bulk file
output_dir = 'bulk'
rows_per_file = 5000  # Total lines (each pair has 2 lines, hence 2500 pairs)

# Run the function
split_bulk_json(input_file, output_dir, rows_per_file)
