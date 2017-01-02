# Basecoat: Formula Search

Basecoat is a CRUD formula indexing tool meant to record formulas for certain colors and store them for future reference.

![Basecoat](http://i.imgur.com/HCaE6dY.png) ![Basecoat Formula](http://i.imgur.com/4HjYcdv.png)

### Prerequisites

* Python 3
* Ubuntu/OSX

## Getting Started

These instructions will get you a copy of the project up and running on your local machine for development and testing purposes.

* Clone from repository
```
git clone git@github.com/cje3295/basecoat.git
```

* Create virtual environment
```
cd basecoat
python3 -m venv venv
source ./venv/bin/activate
```

* Install required python packages
```
pip install -r requirements.txt
```

* Create application configuration file
```
cp example_config.py config.py
```

* Create secret key for config.py
```
python
import os
os.urandom(24).encode('hex')
```
Paste the result in to your configuration file under secret_key

* You can look at the flask script application management options by typing
from the application directory
```
python manage.py --help
```

* Create the database
```
python manage.py db upgrade
```

* Populate the database with fake data for testing
```
python manage.py db populate_db 10
```
The number 10 in this case is the number of entries you want to populate into the database

* Run the local webserver
```
python manage.py runserver -h localhost -p 5000
```

## Managing the database
Through flask-script and the manage.py file you can perform actions on the database to aid in development

#### View Commands
At anytime you can view what commands are possible in manage.py
```
python manage.py db --help
```

#### Populate the database
```
python manage.py db populate_db 10
```
Where the number 10 is the number of entries you want to populate into the database

#### Empty the database
```
python manage.py db empty_db
```

#### Initialize database and create tables
```
python manage.py db upgrade
```

#### Create a new migration
Useful for when you change the database schema in models.py file
```
python manage.py db migrate
python manage.py db upgrade
```

## Built With

* [Flask](http://flask.pocoo.org/) - The web framework used
* [Flask Script](https://flask-script.readthedocs.io/en/latest/) - Application management
* [Flask Migrate](https://flask-migrate.readthedocs.io/en/latest/) - Database migration support/tools

## Authors

* **Clint Edwards** - [Github](https://github.com/cje3295)

## License

This project is licensed under the MIT License - see the [LICENSE.md](LICENSE.md) file for details
