
//When the edit button is clicked we need to open a modal with all the information
//of the row clicked

//A form needs to be posted updating anything that was edited by the modal

//When a row is clicked it should drop down revealing all the information included with the row

//Clicking the plus symbol should open up a modal allowing you to input a new formula


function filterTable(search_string){

    // Overwrite contain selector so we can search case insensitively
    jQuery.expr[':'].contains = function(a, i, m) {
     return jQuery(a).text().toUpperCase().indexOf(m[3].toUpperCase()) >= 0;
    };

    $(".formula").hide();
    $(".formula:contains('" + search_string + "')").show();

    console.log(search_string);
}


function populateViewFormulaModal(formulaID){
    $("#view_modal").find('.modal-title').text('Formula UID: ' + formulaID);
    $("#view_modal").find('.modal-body').load("/formula/" + formulaID);
}


function populateAddFormulaModal(){
    $("#add_modal").find('.modal-title').text('Add New Formula');
    $("#add_modal").find('.modal-body').load("/formula/add");
}


function addExtraBaseForm(){
    $('.base_list .form-inline:first-child').clone().appendTo('.base_list');
}


function addExtraColorantForm(){
    $('.colorant_list .form-inline:first-child').clone().appendTo('.colorant_list');
}


$( document ).ready(function() {

    //Wait for user input on the search bar and then filter table
    $('#search_bar').bindWithDelay("keyup", function() {
        filterTable($(this).val());
    }, 625);

    //Have modal load relevant information on click
    $('.formula').click(function() {
        populateViewFormulaModal($(this).attr('data-formula-id'));
    });

    //Populate add modal
    $('#add_button').click(function() {
        populateAddFormulaModal();
    });

    //Click to add another base edit or add formula page
    $('#add_base').click(function(event) {
        addExtraBaseForm();
        event.stopPropagation();
        event.preventDefault();
    });

    //Click to add another colorant to edit or add formula page
    $('#add_colorant').click(function(event) {
        addExtraColorantForm();
        event.stopPropagation();
        event.preventDefault();
    });

});
