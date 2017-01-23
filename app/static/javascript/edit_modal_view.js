
function addExtraColorantForm(){
    var colorant_list = $(event.target).parent().next('.colorant_list');
    var next_form_id = colorant_list.find('.form-inline').length + 1;

    var colorant_form_item = `
        <div class='form-inline'>
            <div class='form-group'>
                <label for='colorant_product_name${next_form_id}'>Colorant:&nbsp;</label>
                <input type='text' class='form-control' id='colorant_product_name${next_form_id}' name='colorant_product_name${next_form_id}' placeholder='Colorant Product Name' />
            </div>
            <div class='form-group'>
                <input type='text' class='form-control' id='colorant_name${next_form_id}' name='colorant_name${next_form_id}' placeholder='Colorant Name' />
            </div>
            <div class='form-group'>
                <input type='text' class='form-control' id='colorant_amount${next_form_id}' name='colorant_amount${next_form_id}' placeholder='Amount' />
            </div>
        </div>`;

    colorant_list.append(colorant_form_item);
}


$( document ).ready(function() {
    //Click to add another colorant to edit or add formula page
    $('#add_colorant_edit_form').click(function(event) {
        addExtraColorantForm(event);
        event.stopPropagation();
        event.preventDefault();
    });

});
