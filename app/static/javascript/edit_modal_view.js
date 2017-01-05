
function addExtraBaseForm(event){
    var base_list = $(event.target).parent().next('.base_list');
    base_list.find('.form-inline:first-child').clone().appendTo(base_list);
}


function addExtraColorantForm(){
    var colorant_list = $(event.target).parent().next('.colorant_list');
    colorant_list.find('.form-inline:first-child').clone().appendTo(colorant_list);
}


$( document ).ready(function() {
    //Click to add another base edit or add formula page
    $('#add_base_edit_form').click(function(event) {
        addExtraBaseForm(event);
        event.stopPropagation();
        event.preventDefault();
    });

    //Click to add another colorant to edit or add formula page
    $('#add_colorant_edit_form').click(function(event) {
        addExtraColorantForm(event);
        event.stopPropagation();
        event.preventDefault();
    });

});
