import random
import math

class Location:

    def __init__(self):
        x = random.randint(0,10)
        y = random.randint(0,10)
        self.coordinate = (x,y)

    def display_info(self):
        pass

    def distance(self, point):
        dist = math.sqrt((point[0]-self.coordinate[0])**2 + (point[1]-self.coordinate[1])**2)
        return dist

class Consumer(Location):

    def __init__(self, name, student_id):
        super().__init__()
        self.name = name
        self.student_id = student_id

    def display_info(self):
        print("information: ", self.name, self.student_id)
        print("location: ", self.coordinate)

class Pharmacy(Location):

    def __init__(self):
        super().__init__()
        self.quantity = random.randint(0,1000)
    
    def display_info(self):
        print("remain quantity of the pharmacy: ", self.quantity)
        print("location: ", self.coordinate)


def Program(consumer, pharmacy):
    consumer.display_info()
    pharmacy.display_info()

    print("Distance from customer to pharmacy: ", pharmacy.distance(consumer.coordinate))


if __name__ == "__main__":
    customer = Consumer("a", "12312312")
    pharmacy = Pharmacy()
    Program(consumer=customer, pharmacy=pharmacy)