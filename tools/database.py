from flask_migrate import MigrateCommand
from app import db, models

@MigrateCommand.command
def populate_db():
    "Populates database with fake data. Useful for developement"

    dev_formula = models.Formula(color_name='flaky magenta',
                                 color_number='FM-4563',
                                 customer_name='Edwards Inc.')

    dev_colorant = models.Colorant(formula_id=1,
                                   colorant_name='sunset orange',
                                   amount='3')

    dev_base = models.Base(formula_id=1,
                           base_name='simple white',
                           product_name='benjamin moore & co')

    try:
        db.session.add_all([dev_formula, dev_colorant, dev_base])
        db.session.commit()
    except:
        db.session.rollback()
        raise


@MigrateCommand.command
def empty_db():
    """Clears all information from database"""
    meta = db.metadata
    for table in reversed(meta.sorted_tables):
        print 'Cleared table: {}'.format(table)
        db.session.execute(table.delete())
    db.session.commit()
