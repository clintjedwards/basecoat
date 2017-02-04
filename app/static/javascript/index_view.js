
function filterTable(search_string){

    // Overwrite contain selector so we can search case insensitively
    jQuery.expr[':'].contains = function(a, i, m) {
     return jQuery(a).text().toUpperCase().indexOf(m[3].toUpperCase()) >= 0;
    };

    $(".formula").hide();
    $(".formula:contains('" + search_string + "')").show();

}

function convertFormToJSON(form){

    var colorants = $(form).find('.colorant_list .form-inline');
    var colorant_list = [];

    $.each(colorants, function () {
        colorant_product_name = $(this).find('input[id*="colorant_product_name"]').val();
        colorant_name = $(this).find('input[id*="colorant_name"]').val();
        colorant_amount = $(this).find('input[id*="colorant_amount"]').val();

        if (colorant_name){
            colorant_list.push([colorant_product_name, colorant_name, colorant_amount]);
        }
    });


    var form_data = {
                        "formula_id": $(form).find('#formula_id').val(),
                        "formula_name": $(form).find('#formula_name').val(),
                        "formula_number": $(form).find('#formula_number').val(),
                        "base": {"base_product_name": $(form).find('#base_product_name').val(),
                                 "base_name": $(form).find('#base_name').val(),
                                 "base_size": $(form).find('#base_size').val() },
                        "colorant_list": colorant_list,
                        "customer_name": $(form).find('#customer_name').val(),
                        "job_address": $(form).find('#job_address').val(),
                        "notes": $(form).find('#notes').val(),
                    };

    return form_data;

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


function verifyFormulaHasName(form){

    if ($(form).find('#formula_name').val()){
        return true;
    }else{
        return false;
    }

}


$( document ).ready(function() {

    $('[data-toggle=confirmation]').confirmation({
      rootSelector: '[data-toggle=confirmation]',
    });

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

    //Delete button functionality
    $('#view_delete_button').click(function(){
        formulaID = $("#view_modal").find('.modal-title').attr("data-formula-id");

        $.ajax({
            url: '/formula/delete/' + formulaID,
            type: 'DELETE',
            success: function() {
                $('#view_modal #view_close_button').click();
                location.reload(true);
            }
        });
    });

    //Save button functionality
    $('#view_save_button').click(function(){
        if (verifyFormulaHasName($("#edit_form"))){
            $.ajax({
                type: 'POST',
                url: '/formula/add',
                data: JSON.stringify(convertFormToJSON($("#edit_form"))),
                contentType: "application/json",
                dataType: 'json',
                success: function() {
                    $('#view_modal #view_close_button').click();
                    location.reload(true);
                }
            });
        } else {
            $('.top-center').notify({
                message: { text: 'Error: Formula has no name.' },
                type: 'warning',
                closable: false
              }).show();
        }
    });

    //Add formula functionality
    $('#add_save_button').click(function(){
        if (verifyFormulaHasName($("#add_form"))){
            $.ajax({
                type: 'POST',
                url: '/formula/add',
                data: JSON.stringify(convertFormToJSON($("#add_form"))),
                contentType: "application/json",
                dataType: 'json',
                success: function() {
                    $('#add_modal #add_close_button').click();
                    location.reload(true);
                }
            });
        } else {
            $('.top-right').notify({
                message: { text: 'Error: Formula has no name.' },
                type: 'warning',
                closable: false
              }).show();
        }
    });

});
