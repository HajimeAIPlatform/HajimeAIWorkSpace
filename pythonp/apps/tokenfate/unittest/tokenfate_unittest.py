import unittest
from pythonp.apps.tokenfate.src.binance.utils import get_binance_client

class TestTokenfate(unittest.TestCase):

    def test_smoketest(self):
        bc = get_binance_client()
        self.assertTrue(bc)


if __name__ == "__main__":
    unittest.main()