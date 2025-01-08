import unittest
import os

from pythonp.apps.tokenfate.models.db_ops import FortunesDatabase

class TestFortunesDatabase(unittest.TestCase):
    def setUp(self):
        self.test_db_name = 'test_fortunes.db'
        self.db = FortunesDatabase(self.test_db_name)

    def tearDown(self):
        self.db.close()
        if os.path.exists(self.test_db_name):
            os.remove(self.test_db_name)

    def test_connection(self):
        self.assertIsNotNone(self.db.conn)
        self.assertIsNotNone(self.db.cursor)

    def test_insert_and_read_data(self):
        test_data = {
            'sign_url': 'test_url',
            'sign_title': 'Test Sign',
            'sign_unscramble': 'Test Unscramble',
            'sign_overview': 'Test Overview',
            'sign_image_url': 'test_image_url',
            'sign_category': 'Test Category',
            'sign_number': 1
        }
        self.db.insert_or_update_data(test_data)
        
        result = self.db.read_data('Test Category')
        self.assertEqual(len(result), 1)
        self.assertEqual(result[0]['sign_url'], 'test_url')
        self.assertEqual(result[0]['sign_title'], 'Test Sign')

    def test_update_data(self):
        test_data = {
            'sign_url': 'update_test_url',
            'sign_title': 'Update Test Sign',
            'sign_unscramble': 'Update Test Unscramble',
            'sign_overview': 'Update Test Overview',
            'sign_image_url': 'update_test_image_url',
            'sign_category': 'Update Test Category',
            'sign_number': 2
        }
        self.db.insert_or_update_data(test_data)
        
        update_data = {
            'sign_title': 'Updated Sign',
            'sign_overview': 'Updated Overview'
        }
        self.db.update_data('update_test_url', update_data)
        
        result = self.db.read_data('Update Test Category')
        self.assertEqual(result[0]['sign_title'], 'Updated Sign')
        self.assertEqual(result[0]['sign_overview'], 'Updated Overview')

    def test_delete_data(self):
        test_data = {
            'sign_url': 'delete_test_url',
            'sign_title': 'Delete Test Sign',
            'sign_unscramble': 'Delete Test Unscramble',
            'sign_overview': 'Delete Test Overview',
            'sign_image_url': 'delete_test_image_url',
            'sign_category': 'Delete Test Category',
            'sign_number': 3
        }
        self.db.insert_or_update_data(test_data)
        
        self.db.delete_data('delete_test_url')
        
        result = self.db.read_data('Delete Test Category')
        self.assertEqual(len(result), 0)

    def test_randomly_choose_sign_by_weight(self):
        for i in range(5):  # Insert 5 test signs
            test_data = {
                'sign_url': f'random_test_url_{i}',
                'sign_title': f'Random Test Sign {i}',
                'sign_unscramble': f'Random Test Unscramble {i}',
                'sign_overview': f'Random Test Overview {i}',
                'sign_image_url': f'random_test_image_url_{i}',
                'sign_category': 'Random Test Category',
                'sign_number': i
            }
            self.db.insert_or_update_data(test_data)
        
        chosen_sign = self.db.randomly_choose_sign_by_weight()
        self.assertIsNotNone(chosen_sign)
        self.assertIn('sign_overview', chosen_sign)
        self.assertIn('sign_image_url', chosen_sign)
        self.assertIn('sign_category', chosen_sign)

if __name__ == '__main__':
    unittest.main()