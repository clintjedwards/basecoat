
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
    });

    //Revert edit button swapout
    $('#view_close_button').click(function(){
        $('#view_edit_button').show();
        $('#view_save_button').hide();
    });

    //Save button functionality
    $('#view_save_button').click(function(){
        console.log('We sent a post request here to edit a formula');
    });

    //Add formula functionality
    $('#add_save_button').click(function(){
        console.log($("#add_form").serialize());
        $.post("/formula/add", $("#add_form").serialize());
    });

});
