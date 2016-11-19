

class Formula(object):

    def __init__(self, color_name="", color_number=0):
        self.color_name = color_name
        self.color_number = color_number
        self.customer_name = ""
        self.colorants = {}
        self.base_to_product_map = {}

    def __str__(self):
        return str(self.__dict__)


exampleform = Formula()
print exampleform
