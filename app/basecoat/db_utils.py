from app import db, models

def insert_into_db(table_model_class, **kwargs):
    """ Insert information into database

    Returns the database object and then assigns values by kwarguments.

    Args:
        table_model_class: The name of the model of the database
        kwargs: Database column names and value
    """

    table_model = getattr(models, table_model_class)
    query = table_model(**kwargs)
    db.session.add(query)
    db.session.commit()

def update_db(table_model_class, key_name, key_value, **kwargs):
    """ Update database object

    Grabs single database object by column name and column value
    and then updates the object based on key word arguments.

    Args:
        table_model_class: The name of the model of the database
        key_name: Column name. ex Users
        key_value: Column value. ex clint.edwards
        kwargs: Database column names and value
    """

    table_model = getattr(models, table_model_class)
    key = getattr(table_model, key_name)
    table_info = table_model.query.filter(key == key_value).first()

    for argument in kwargs:
        setattr(table_info, argument, kwargs[argument])

    db.session.commit()

def get_table(table_model_class):
    """ Return database table

    Args:
        table_model_class: The name of the model of the database
    """

    table_model = getattr(models, table_model_class)
    table_info = table_model.query.all()

    return table_info

def get_object_from_table(table_model_class, key_name, key_value):
    """ Return database object

    Returns single database object based on column name and value

    Args:
        table_model_class: The name of the model of the database
        key_name: Column name. ex Users
        key_value: Column value. ex clint.edwards
    """

    table_model = getattr(models, table_model_class)
    key = getattr(table_model, key_name)
    table_info = table_model.query.filter(key == key_value).first()

    return table_info

def delete_from_db(table_model_class, key_name, key_value):
    """ Delete database object

    Return database object by column name and value and delete it.

    Args:
        table_model_class: The name of the model of the database
        key_name: Column name. ex Users
        key_value: Column value. ex clint.edwards
    """

    table_model = getattr(models, table_model_class)
    key = getattr(table_model, key_name)
    table_model.query.filter(key == key_value).delete()
    db.session.commit()
