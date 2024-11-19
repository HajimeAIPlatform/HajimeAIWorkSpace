import json
import sqlite3
import os
import random
from typing import List, Dict, Any
from pythonp.apps.tokenfate.models.weight_random import WeightRandom
import os

data_source = 'pythonp/apps/tokenfate/static/assets/refined_lots.json'

class FortunesDatabase:
    def __init__(self, db_name: str = 'fortunes.db'):
        self.db_name = db_name
        self.conn = None
        self.cursor = None
        self.rm_db()
        self.connect()
        self.create_table()
        # json_data = self.load_json_data(data_source)
        json_data = self.load_data(data_source)
        for item in json_data:
            self.insert_or_update_data(item)
        print("Data loaded.")

    def rm_db(self):
        if os.path.exists(self.db_name):
            os.remove(self.db_name)
            print(f"Database '{self.db_name}' removed.")
        else:
            print(f"Database '{self.db_name}' does not exist.")

    def connect(self):
        db_exists = os.path.exists(self.db_name)
        if not db_exists:
            print(f'Starting new database "{self.db_name}"...')
            self.conn = sqlite3.connect(self.db_name)
            self.cursor = self.conn.cursor()
            print(f"Database '{self.db_name}' created.")        
        else:
            print(f"Connected to existing database '{self.db_name}'.")

    def create_table(self):
        self.cursor.execute('''
        CREATE TABLE IF NOT EXISTS signs (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            weight INTEGER,
            sign_level TEXT,
            sign_from TEXT,
            sign_text TEXT
        )
        ''')
        self.conn.commit()
        print("Table 'signs' is ready.")

    def insert_or_update_data(self, data: Dict[str, Any]):
        weight = random.randint(0, 1000)  # Generate random weight
        self.cursor.execute('''
        INSERT OR REPLACE INTO signs (
            weight, sign_level, sign_from, sign_text
        ) VALUES (?, ?, ?, ?)
        ''', (
            weight, data['sign_level'], data['sign_from'], data['sign_text']
        ))
        self.conn.commit()

    def read_data(self, category: str = None) -> List[Dict[str, Any]]:
        if category:
            self.cursor.execute('SELECT * FROM signs WHERE sign_category = ?', (category,))
        else:
            self.cursor.execute('SELECT * FROM signs')
        
        columns = [column[0] for column in self.cursor.description]
        return [dict(zip(columns, row)) for row in self.cursor.fetchall()]
    
    def get_lot_by_id(self, id: int) -> Dict[str, Any]:
        self.cursor.execute('SELECT * FROM signs WHERE id = ?', (id,))
        return dict(zip([column[0] for column in self.cursor.description], self.cursor.fetchone()))
    
    def read_id_and_weight(self):
        sql = "SELECT id, weight FROM signs"
        self.cursor.execute(sql)
        return [(row[0], row[1]) for row in self.cursor.fetchall()]

    # 根据权重随机抽签
    def randomly_choose_sign_by_weight(self) -> Dict[str, Any]:
        id_weight_pairs = self.read_id_and_weight()
        random_weights = WeightRandom(id_weight_pairs)
        id = random_weights.choose()
        return self.get_lot_by_id(id)
        # sql = 'SELECT sign_level, sign_from, sign_text FROM signs WHERE id = ?'
        # self.cursor.execute(sql, (id,))
        # return [dict(zip(['id', 'sign_level', 'sign_from', 'sign_text'], row)) for row in self.cursor.fetchall()][0]

    def update_data(self, sign_url: str, data: Dict[str, Any]):
        if 'weight' not in data:
            data['weight'] = random.randint(0, 1000)  # Generate new random weight if not provided
        update_fields = ', '.join([f"{key} = ?" for key in data.keys()])
        query = f"UPDATE signs SET {update_fields} WHERE sign_url = ?"
        self.cursor.execute(query, list(data.values()) + [sign_url])
        self.conn.commit()

    def delete_data(self, sign_url: str):
        self.cursor.execute('DELETE FROM signs WHERE sign_url = ?', (sign_url,))
        self.conn.commit()

    def close(self):
        if self.conn:
            self.conn.close()
            print(f"Connection to '{self.db_name}' closed.")

    # def load_json_data(self, file_path: str) -> List[Dict[str, Any]]:
    #     with open(file_path, 'r', encoding='utf-8') as file:
    #         raws = json.load(file)
    #         all_signs = []
    #         for raw in raws:
    #             all_signs.extend(raw.get("related_signs", []))
    #         return all_signs
    def load_data(self, file_path: str) -> List[Dict[str, Any]]:
        with open(file_path, 'r', encoding='utf-8') as file:
            signs = json.load(file)
            return signs

def sync_database_with_json(db: FortunesDatabase, json_data: List[Dict[str, Any]]):
    # Get existing URLs in the database
    existing_urls = set(item['sign_url'] for item in db.read_data())

    # Update or insert data from JSON
    for item in json_data:
        db.insert_or_update_data(item)

    # Delete items that are in the database but not in the JSON
    json_urls = set(item['sign_url'] for item in json_data)
    for url in existing_urls - json_urls:
        db.delete_data(url)

def main():
    # Initialize database
    db = FortunesDatabase()
    # selection = db.randomly_choose_sign_by_weight()
    # print(selection)
    
    # Close database connection
    db.close()

if __name__ == "__main__":
    main()
