// DOM Ready =============================================================
$(document).ready(function() {

    // Populate the vehicle table on initial page load
    populateTable();

    // Mileage link click
    $('#vehicleList table tbody').on('click', 'td a.linkshowmileage', showMileageInfo);

    // Add Vehicle button click
    $('#btnAddVehicle').on('click', addVehicle);

    // Add Mileage button click
    $('#btnAddMileage').on('click', addMileage);

    // Show History button click
    $('#btnShowMileageLog').on('click', showMileageLog);

});

// Functions =============================================================

// Fill table with data
function populateTable() {

    // Empty content string
    var tableContent = '';

    // jQuery AJAX call for JSON
    $.getJSON( '/vehicles', function( data ) {

        // For each item in our JSON, add a table row and cells to the content string
        $.each(data, function(){
            tableContent += '<tr>';
            tableContent += '<td><a href="/vehicles/mileage/' + this.Key + '" class="linkshowmileage" rel="' + this.Key + '">' + this.Key + '</a></td>';
            tableContent += '<td>' + this.Record.type + '</td>';
            tableContent += '<td>' + this.Record.manufactor + '</td>';
            tableContent += '<td>' + this.Record.model + '</td>';
            tableContent += '<td>' + this.Record.Registration + '</td>';
            tableContent += '</tr>';
        });

        // Inject the whole content string into our existing HTML table
        $('#vehicleList table tbody').html(tableContent);
    });
};

// Show Mileage Info
function showMileageInfo(event) {

    // Prevent Link from Firing
    event.preventDefault();

    // Retrieve id from link rel attribute
    var thisId = $(this).attr('rel');

    // jQuery AJAX call for JSON
    $.getJSON( '/vehicles/mileage/' + thisId, function( data ) {
        //Populate Info Box
        $('#mileageInfoDate').text(data.date);
        $('#mileageInfoValue').text(data.mileage);
        $('#mileageInfoComment').text(data.comment);
    }).fail( function () {
        //Clean up Info Box
        $('#mileageInfoDate').text('');
        $('#mileageInfoValue').text('');
        $('#mileageInfoComment').text('');
    });
    $('#mileageInfoId').text(thisId);
    $('#btnAddMileage').removeAttr('disabled');
    $('#btnShowMileageLog').removeAttr('disabled');
};

// Add Vehicle
function addVehicle(event) {
    event.preventDefault();

    // Super basic validation - increase errorCount variable if any fields are blank
    var errorCount = 0;
    $('#addVehicle input').each(function(index, val) {
        if($(this).val() === '') { errorCount++; }
    });

    // Check and make sure errorCount's still at zero
    if(errorCount === 0) {

        // If it is, compile all vehicle info into one object
        var vehicle = {
            'id': $('#addVehicle fieldset input#inputVehicleId').val(),
            'type': $('#addVehicle fieldset input#inputVehicleType').val(),
            'manufactor': $('#addVehicle fieldset input#inputVehicleManufactor').val(),
            'model': $('#addVehicle fieldset input#inputVehicleModel').val(),
            'registration': $('#addVehicle fieldset input#inputVehicleRegistration').val()
        }

        // Use AJAX to post the object to our adduser service
        $.ajax({
            type: 'POST',
            data: vehicle,
            url: '/vehicles/addvehicle',
            dataType: 'JSON'
        }).done(function( response ) {

            // Check for successful (blank) response
            if (response.msg === '') {

                // Clear the form inputs
                $('#addVehicle fieldset input').val('');

                // Update the table
                populateTable();

            }
            else {

                // If something goes wrong, alert the error message that our service returned
                alert('Error: ' + response.msg);

            }
        });
    }
    else {
        // If errorCount is more than 0, error out
        alert('Please fill in all fields');
        return false;
    }
};

// Add Vehicle
function addMileage(event) {
    event.preventDefault();

    // Super basic validation - increase errorCount variable if any fields are blank
    var errorCount = 0;
    $('#addMileage input').each(function(index, val) {
        if($(this).val() === '') { errorCount++; }
    });

    // Check and make sure errorCount's still at zero
    if(errorCount === 0) {

        var id = $('#mileageInfoId').text();
        // If it is, compile all mileage info into one object
        var mileage = {
            'date': $('#addMileage fieldset input#inputMileageDate').val(),
            'mileage': $('#addMileage fieldset input#inputMileageValue').val(),
            'comment': $('#addMileage fieldset input#inputMileageComment').val(),
        }

        // Use AJAX to post the object to our adduser service
        $.ajax({
            type: 'POST',
            data: mileage,
            url: '/vehicles/addmileage/' + id,
            dataType: 'JSON'
        }).done(function( response ) {

            // Check for successful (blank) response
            if (response.msg === '') {

                // Clear the form inputs
                $('#addMileage fieldset input').val('');

                // Update the table
                $('#vehicleList a[rel='+id+']').trigger('click');

            }
            else {

                // If something goes wrong, alert the error message that our service returned
                alert('Error: ' + response.msg);

            }
        });
    }
    else {
        // If errorCount is more than 0, error out
        alert('Please fill in all fields');
        return false;
    }
};

// Show Mileage Info
function showMileageLog(event) {

    // Prevent Link from Firing
    event.preventDefault();

    // Retrieve id from link rel attribute
    var thisId = $('#mileageInfoId').text();

    // jQuery AJAX call for JSON
    $.getJSON( '/vehicles/mileageLog/' + thisId, function( data ) {

        // Empty content string
        var tableContent = '';

        // For each item in our JSON, add a table row and cells to the content string
        $.each(data, function(){
            tableContent += '<tr>';
            tableContent += '<td>' + this.date + '</td>';
            tableContent += '<td>' + this.mileage + '</td>';
            tableContent += '</tr>';
        });

        // Inject the whole content string into our existing HTML table
        $('#mileageLog table tbody').html(tableContent);
    });
};

