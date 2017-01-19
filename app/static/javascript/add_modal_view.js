
function addExtraBaseForm(event){
    var base_list = $(event.target).parent().next('.base_list');
    var next_form_id = base_list.find('.form-inline').length + 1;

    var base_form_item = `
        <div class='form-inline'>
            <div class='form-group'>
                <label for='InputBase${next_form_id}'>Base:&nbsp;</label>
                <input type='text' class='form-control' id='InputBase${next_form_id}' name='InputBase${next_form_id}' placeholder='Base Name' />
            </div>
            <div class='form-group'>
                <label class='sr-only' for='InputProductName${next_form_id}'>Base Product Name</label>
                <input type='text' class='form-control' id='InputProductName${next_form_id}' name='InputProductName${next_form_id}' placeholder='Base Product Name' />
            </div>
        </div>`;

    base_list.append(base_form_item);
}


function addExtraColorantForm(){
    var colorant_list = $(event.target).parent().next('.colorant_list');
    var next_form_id = colorant_list.find('.form-inline').length + 1;

    var colorant_form_item = `
        <div class='form-inline'>
            <div class='form-group'>
                <label for='InputColorant${next_form_id}'>Colorant:&nbsp;</label>
                <input type='text' class='form-control' id='InputColorant${next_form_id}' name='InputColorant${next_form_id}' placeholder='Colorant Name' />
            </div>
            <div class='form-group'>
                <label class='sr-only' for='InputAmount${next_form_id}'>Amount</label>
                <input type='number' class='form-control' id='InputAmount${next_form_id}' name='InputAmount${next_form_id}' placeholder='Amount (in gallons)' />
            </div>
        </div>`;

    colorant_list.append(colorant_form_item);
}


$( document ).ready(function() {
    //Click to add another base edit or add formula page
    $('#add_base_add_form').click(function(event) {
        addExtraBaseForm(event);
        event.stopPropagation();
        event.preventDefault();
    });

    //Click to add another colorant to edit or add formula page
    $('#add_colorant_add_form').click(function(event) {
        addExtraColorantForm(event);
        event.stopPropagation();
        event.preventDefault();
    });

});
