
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

}

function convertFormToJSON(form){

    var bases = $('.base_list .form-inline');
    var colorants = $('.colorant_list .form-inline');

    var base_dict = {};
    var colorant_dict = {};

    $.each(bases, function () {
       base_name = $(this).find('input[id*="InputBase"]').val();
       base_product_name = $(this).find('input[id*="InputProductName"]').val();

       if (base_name){
           base_dict[base_name] = base_product_name
       };
    });

    $.each(colorants, function () {
       colorant_name = $(this).find('input[id*="InputColorant"]').val();
       colorant_amount = $(this).find('input[id*="InputAmount"]').val();

       if (colorant_name){
           colorant_dict[colorant_name] = colorant_amount
       };
    });

    var form_data = {
                        "formula_id": $('#InputFormulaID').val(),
                        "formula_name": $('#InputFormulaName').val(),
                        "formula_number": $('#InputFormulaNumber').val(),
                        "base_list": base_dict,
                        "colorant_list": colorant_dict,
                        "customer": $('#InputCustomer').val(),
                        "summary": $('#InputSummary').val(),
                        "notes": $('#InputNotes').val(),
                    };

}


function populateViewFormulaModal(formulaID){
    $("#view_modal").find('.modal-title').attr("data-formula-id", formulaID);
    $("#view_modal").find('.modal-title').text('Formula UID: ' + formulaID);
    $("#view_modal").find('.modal-body').load("/formula/" + formulaID);
}


function populateAddFormulaModal(){
    $("#add_modal").find('.modal-title').text('Add New Formula');
    $("#add_modal").find('.modal-body').load("/formula/add");
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

    //Edit button functionality
    $('#view_edit_button').click(function() {
        formulaID = $("#view_modal").find('.modal-title').attr("data-formula-id");
        $("#view_modal").find('.modal-body').load("/formula/edit/" + formulaID);

        $('#view_edit_button').hide();
        $('#view_save_button').show();
        $('#view_delete_button').show();
    });

    //Revert edit button swapout
    $('#view_close_button').click(function(){
        $('#view_edit_button').show();
        $('#view_save_button').hide();
        $('#view_delete_button').hide();
    });

    //Save button functionality
    $('#view_save_button').click(function(){
        convertFormToJSON();
        //$.post("/formula/add", $("#edit_form").serialize());
    });

    //Add formula functionality
    $('#add_save_button').click(function(){
        convertFormToJSON();
        //$.post("/formula/add", $("#add_form").serialize());
    });

});
