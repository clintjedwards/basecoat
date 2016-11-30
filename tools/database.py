from flask_migrate import MigrateCommand
from app import db, models

@MigrateCommand.command
def populate_db():
    "Populates database with fake data. Useful for developement"

    dev_formula1 = models.Formula(color_name='Soft Red',
                                  color_number='SR-EC644B',
                                  customer_name='Edwards Inc.',
                                  summary='Etiam iaculis nulla ac ex euismod tempor. Sed suscipit lorem id urna porttitor ullamcorper.',
                                  notes='In egestas magna eu turpis condimentum venenatis.')
    dev_formula2 = models.Formula(color_name='Ripe Lemon',
                                  color_number='RL-F7CA18',
                                  customer_name='Edwards Inc.',
                                  summary='Curabitur facilisis diam a pharetra scelerisque. In quam mi, dapibus at ligula et, blandit ultricies dolor.',
                                  notes='Aliquam eget urna euismod, consectetur tellus et, auctor metus.')
    dev_formula3 = models.Formula(color_name='Light Wisteria',
                                  color_number='LW-BE90D4',
                                  customer_name='Edwards Inc.',
                                  summary='',
                                  notes='')
    dev_formula4 = models.Formula(color_name='Testcolorlongname',
                                  color_number='TC-BE90D4-BE90D4-BE90D4',
                                  customer_name='Edwards Inc.',
                                  summary='Curabitur facilisis diam a pharetra scelerisque. In quam mi, dapibus at ligula et, blandit ultricies dolor.' +
                                          'Curabitur facilisis diam a pharetra scelerisque. In quam mi, dapibus at ligula et, blandit ultricies dolor.' +
                                          'Curabitur facilisis diam a pharetra scelerisque. In quam mi, dapibus at ligula et, blandit ultricies dolor.',
                                  notes='')

    dev_formula5 = models.Formula(color_name='Tcs',
                                  color_number='TC',
                                  customer_name='Edwards Inc.',
                                  summary='',
                                  notes='')

    dev_colorant = models.Colorant(formula_id=1,
                                   colorant_name='sunset orange',
                                   amount='3')

    dev_base = models.Base(formula_id=1,
                           base_name='simple white',
                           product_name='benjamin moore & co')

    try:
        db.session.add_all([dev_formula1, dev_formula2, dev_formula3, dev_formula4, dev_formula5, dev_colorant, dev_base])
        db.session.commit()
        print 'Populated database'
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
