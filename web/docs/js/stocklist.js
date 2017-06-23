function loadStockList(symbol,tabId) {
  
        $('#statustext').html('sending..');
        
        $.ajax({
            url: '/stocks?symbol='+symbol,
            type: 'post',
            dataType: 'json',
            success: function (data) {
                $('#'+tabId).empty()
                $.each(data, function(index, e) {
                        tr = "<tr><td>" + e.symbol + "</td><td>" + e.Price + "</td><td>" + e.Quantity + "</td><td>" + e.InsertDate + "</td><td>" + e.Operation + "</td>"
                        tr+= "<td><button id='"+e.id +"' class='delstock btn-danger btn-xs'>Delete</button></td></tr>"
                        //alert(tr)
                        $('#'+tabId).append(tr)
                });
                registerDelButtons()
                $('#stocks').DataTable();
                $('#statustext').html('done');
            },
            error: function(data) {
                $('#statustext').html('failed');
            },
            //data: person
        });
}

function refreshTable(stockTable    ){
    stockTable.ajax.reload();
}

function registerDelButtons(){
    $(".delstock").unbind();
    $(".delstock").bind("click",function(){     
            DeleteStock($(this).attr('id'))
        });
}

function DeleteStock(delId){  
    var req = {
            id: delId
        }
        callBackEnd('/stocks/del',req)
}

function callBackEnd(url,reqData){
        $('#statustext').html('sending..');
        $.ajax({
            url: url,
            type: 'post',
            contentType: 'application/json',
            dataType: 'json',
            success: function (data) {          
                $('#statustext').html('done');
            },
            error: function(data) {
                $('#statustext').html('failed');
            },
            data: JSON.stringify(reqData)
        });
}

function buyCall() {
    var req = {
            symbol: $("#tick").val(),
            Price:parseFloat($("#price").val()),
            Quantity:parseInt($("#quant").val()),
            Operation: "Buy",
            InsertDate: $("#datepicker").val(),
            Tax: parseFloat($("#tax").val())
        }
       callBackEnd('/stocks/add',req)
}

function sellCall() {
    var req = {
            symbol: $("#tick").val(),
            Price:parseFloat($("#price").val()),
            Quantity:parseInt($("#quant").val()),
            Operation: "Sell",
            InsertDate: $("#datepicker").val()
        }
       callBackEnd('/stocks/add',req)
}

function EnableTableSorter() {  
  $("#stocks").tablesorter({sortList: [[2,0]]}); 
}

function EnableDatePicker() {
      $("#datepicker").datepicker({
             changeMonth: true,//this option for allowing user to select month
             changeYear: true, //this option for allowing user to select from year range
             dateFormat: 'dd/mm/yy'
        });
}