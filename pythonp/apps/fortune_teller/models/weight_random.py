import random
import bisect
from typing import List, Any

class WeightRandom:
    def __init__(self, id_weight_pairs: List[tuple[Any, float]]):
        self.ids = []
        self.weight_ranges = []
        self.total_weight = 0
        self.cal_weights(id_weight_pairs)

    def cal_weights(self, id_weight_pairs: List[tuple[Any, float]]):
        self.ids = []
        self.weight_ranges = []
        self.total_weight = 0

        for id, weight in id_weight_pairs:
            if weight <= 0:
                continue

            self.ids.append(id)
            self.total_weight += weight
            self.weight_ranges.append(self.total_weight)

    def choose(self) -> Any:
        random_value = random.uniform(0, self.total_weight)
        index = bisect.bisect_left(self.weight_ranges, random_value)
        return self.ids[index]

def main():
    sample_count = 1_000_000

    id_weight_pairs = [
        ("server1", 1.0),
        ("server2", 3.0),
        ("server3", 2.0)
    ]

    weight_random = WeightRandom(id_weight_pairs)

    statistics = {}

    for _ in range(sample_count):
        choice = weight_random.choose()
        statistics[choice] = statistics.get(choice, 0) + 1

    for k, v in statistics.items():
        hit = v / sample_count
        print(f"{k}, hit: {hit:.4f}")

if __name__ == "__main__":
    main()