from flask import render_template, request, jsonify
from app.basecoat import db_utils

from app import app, db, models


@app.route('/')
def index():
    formula_table = db_utils.get_table('Formula')
    formula_list = [formula for formula in formula_table]
    return render_template('index.html',
                           formula_list=formula_list)


@app.route('/formula/<int:formula_id>')
def get_formula(formula_id):
    formula = db_utils.get_object_from_table('Formula', 'id', formula_id)[0]
    colorant_list = db_utils.get_object_from_table('Colorant', 'formula_id', formula_id)
    base_list = db_utils.get_object_from_table('Base', 'formula_id', formula_id)
    return render_template('view_formula.html',
                           formula=formula,
                           colorant_list=colorant_list,
                           base_list=base_list)


@app.route('/formula/add', methods=['GET', 'POST'])
def add_formula():
    if request.method == 'POST':
        form_data = {key: value.strip() for key, value in request.form.items()}
        if "InputFormulaID" in form_data:
            for thing in form_data:
                print(thing)
                print(form_data[thing])
        else:
            new_formula = models.Formula(color_name=form_data['InputFormulaName'].title(),
                                         color_number=form_data['InputFormulaNumber'].title(),
                                         customer_name=form_data['InputCustomer'].title(),
                                         summary=form_data['InputSummary'].capitalize(),
                                         notes=form_data['InputNotes'].capitalize())
            db.session.add(new_formula)
            db.session.flush()

            print(new_formula.id)
            try:
                db.session.commit()
            except:
                db.session.rollback()
                raise

        return jsonify({'success':True}), 200
    else:
        return render_template('add_formula.html')


@app.route('/formula/edit/<int:formula_id>')
def edit_formula(formula_id):
    formula = db_utils.get_object_from_table('Formula', 'id', formula_id)[0]
    colorant_list = db_utils.get_object_from_table('Colorant', 'formula_id', formula_id)
    base_list = db_utils.get_object_from_table('Base', 'formula_id', formula_id)
    return render_template('edit_formula.html',
                           formula=formula,
                           colorant_list=colorant_list,
                           base_list=base_list)

# dev_formula = models.Formula(color_name=fake.color_name() + " " + fake.safe_color_name(),
#                              color_number=fake.hex_color(),
#                                 customer_name=fake.company(),
#                         summary=fake.text(max_nb_chars=random.randint(50, 200)),
#                     notes=fake.paragraph(nb_sentences=3, variable_nb_sentences=True))
#
# try:
#     db.session.add_all([dev_formula])
#     db.session.commit()
# #.strip()
#.title()
#    color_name = db.Column(db.String(64))
#    color_number = db.Column(db.String(64))
#    customer_name = db.Column(db.String(64))
#    summary = db.Column(db.String(64))
#    notes = db.Column(db.Text())
