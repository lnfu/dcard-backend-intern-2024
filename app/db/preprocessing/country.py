import csv

def process_csv(file_path):
    with open(file_path, 'r', newline='') as file:
        reader = csv.reader(file)
        for row in reader:
            if len(row) >= 2:
                col1 = "'" + row[0] + "'"
                col2 = "'" + row[1] + "'"
                print(f'({col1}, {col2}),')

file_path = 'country.csv'
process_csv(file_path)