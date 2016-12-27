from flask_migrate import MigrateCommand
from elizabeth import Text, Business
from app import db, models
import random

text = Text('en')
business = Business('en')

@MigrateCommand.command
def populate_db(number_of_entries):
    "Populates database with fake data. Useful for developement"

    number_of_entries = int(number_of_entries)

    for entry in range(0, number_of_entries):
        dev_formula = models.Formula(color_name=text.color() + " " + text.color(),
                                     color_number=text.hex_color(),
                                     customer_name=business.company(),
                                     summary=text.sentence(),
                                     notes=text.text(quantity=5))

        try:
            db.session.add_all([dev_formula])
            db.session.commit()
        except:
            db.session.rollback()
            raise

    print('Populated formulas')

    for entry in range(0, number_of_entries):

        dev_colorant = models.Colorant(formula_id=random.randint(1, number_of_entries),
                                       colorant_name=text.color() + " " + text.color(),
                                       amount=random.randint(1, 10))

        try:
            db.session.add_all([dev_colorant])
            db.session.commit()
        except:
            db.session.rollback()
            raise

    print('Populated colorants')

    for entry in range(0, number_of_entries):

        dev_base = models.Base(formula_id=random.randint(1, number_of_entries),
                               base_name=text.color() + " " + text.color(),
                               product_name=business.company())

        try:
            db.session.add_all([dev_base])
            db.session.commit()
        except:
            db.session.rollback()
            raise

    print('Populated bases')


@MigrateCommand.command
def empty_db():
    """Clears all information from database"""
    meta = db.metadata

    for table in reversed(meta.sorted_tables):
        print('Cleared table: {}'.format(table))
        db.session.execute(table.delete())

    db.session.commit()
