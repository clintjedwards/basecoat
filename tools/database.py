from flask_migrate import MigrateCommand
from faker import Faker
from app import db, models
import random

fake = Faker()

@MigrateCommand.command
def populate_db(number_of_entries):
    "Populates database with fake data. Useful for developement"

    number_of_entries = int(number_of_entries)

    for entry in range(0, number_of_entries):
        dev_formula = models.Formula(formula_name=fake.color_name() + " " + fake.safe_color_name(),
                                     formula_number=fake.hex_color(),
                                     customer_name=fake.company(),
                                     summary=fake.text(max_nb_chars=random.randint(50, 200)),
                                     notes=fake.paragraph(nb_sentences=3, variable_nb_sentences=True))

        try:
            db.session.add_all([dev_formula])
            db.session.commit()
        except:
            db.session.rollback()
            raise

    print('Populated formulas')

    for entry in range(0, number_of_entries):

        dev_colorant = models.Colorant(formula_id=random.randint(1, number_of_entries),
                                       colorant_name=fake.color_name() + " " + fake.safe_color_name(),
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
                               base_name=fake.color_name() + " " + fake.safe_color_name(),
                               product_name=fake.company())

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
