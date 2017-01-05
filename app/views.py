from flask import render_template, request, jsonify
from app.basecoat import db_utils as db

from app import app


@app.route('/')
def index():
    formula_table = db.get_table('Formula')
    formula_list = [formula for formula in formula_table]
    return render_template('index.html',
                           formula_list=formula_list)


@app.route('/formula/<int:formula_id>')
def get_formula(formula_id):
    formula = db.get_object_from_table('Formula', 'id', formula_id)[0]
    colorant_list = db.get_object_from_table('Colorant', 'formula_id', formula_id)
    base_list = db.get_object_from_table('Base', 'formula_id', formula_id)
    return render_template('view_formula.html',
                           formula=formula,
                           colorant_list=colorant_list,
                           base_list=base_list)


@app.route('/formula/add', methods=['GET', 'POST'])
def add_formula():
    if request.method == 'POST':
        for thing in request.form:
            print(thing + ": ")
            print(request.form[thing])
        return jsonify({'success':True}), 200
    else:
        return render_template('add_formula.html')


@app.route('/formula/edit/<int:formula_id>')
def edit_formula(formula_id):
    formula = db.get_object_from_table('Formula', 'id', formula_id)[0]
    colorant_list = db.get_object_from_table('Colorant', 'formula_id', formula_id)
    base_list = db.get_object_from_table('Base', 'formula_id', formula_id)
    return render_template('edit_formula.html',
                           formula=formula,
                           colorant_list=colorant_list,
                           base_list=base_list)
