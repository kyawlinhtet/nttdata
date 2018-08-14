var map,marker;
var resetFlag = false;
var paginationKeyword,paginationLocation = "";
jQuery("#page").val(1);
queryYelp();
jQuery("#txtKeyword,#txtLocation").on("change",function(event){
    resetFlag = true;
});
jQuery("#btnSearch").on("click",function(event){
    if(resetFlag){
        resetFlag = false;
        jQuery("#page").val(1);
    }
    queryYelp();
});
$('#modalMap').on('show.bs.modal', function (e) {
    jQuery("#mapTitle").text(jQuery(e.relatedTarget).data("bizname"));
})
$('#modalMap').on('shown.bs.modal', function (e) {
    initMap(jQuery(e.relatedTarget).data("lat"),jQuery(e.relatedTarget).data("lng"));
})
function queryYelp(){
    jQuery(".result-list").html(getLoaderHtml());
    jQuery('.pagination').hide();
    jQuery.get("/get",$('#frmSearch').serialize(), function(data) {
        if(data){   
            data = JSON.parse(data);
            var str = "";
            jQuery.each(data.businesses, function(index, value) {
                str += CreateBizCard(value);
            });
            jQuery(".result-list").html(str);
            jQuery('.pagination').show();
            CreatePagination(data.total)
        }else{
            jQuery(".result-list").html('<div class="mx-auto d-none"><h3>Search No Result</h3></div>');
        }
    })
}
function CreateBizCard(biz){
    var str = '<div class="card">';
    str += '<a target="_blank"  href="'+biz.url+'"><img class="card-img-top" src="'+biz.image_url+'" alt="'+biz.name+'"></a>';
    str += '<div class="card-body">';
    str += '<h4 class="card-title"><a target="_blank"  href="'+biz.url+'">'+biz.name+'</a></h4>';
    str += '<p class="card-text"><i class="fas fa-phone mr-1"></i>'+biz.display_phone+'</p>';
    str += '<p class="card-text"><i class="fas fa-address-book mr-1"></i>'+biz.location.display_address[0]+biz.location.display_address[1]+'</p>';
    str += '<p class="card-text"><i class="fas fa-star mr-1"></i>'+biz.rating+'</p>';
    str += '</div>';
    str += '<div class="card-footer">';
    str += '<div class="text-center"><a href="#modalMap" class="btn" data-toggle="modal" data-bizname="'+biz.name+'" data-lat="'+biz.coordinates.latitude+'" data-lng="'+biz.coordinates.longitude+'"><i class="fas fa-map-marked-alt mr-1"></i>View on Map</a></div>';
    str += '</div></div>';
    return str;
}
function CreatePagination(total){
    var itemPerPage = 15;
    var currentPage = parseInt(jQuery("#page").val());
    paginationKeyword = jQuery("#txtKeyword").val();
    paginationLocation = jQuery("#txtLocation").val();
    if(!currentPage)currentPage=1;
    jQuery('.pagination').twbsPagination({
        totalPages: Math.ceil(total/itemPerPage),
        visiblePages: 5,
        startPage: currentPage,
        first: '&laquo;',
        prev: '<',
        next: '>',
        last: '&raquo;',
        hideOnlyOnePage: true,
        initiateStartPageClick: false,
        onPageClick: function (event, page) {
            jQuery("#txtKeyword").val(paginationKeyword);
            jQuery("#txtLocation").val(paginationLocation);
            jQuery("#page").val(page);
            queryYelp();
        }
    });   
}
// Initialize and add the map
function initMap(lat,lng) {
    // The location of Uluru
    var coordinates = {lat: 1.370270, lng: 103.851959};
    var zoom = 11;
    if(lat && lng){
        coordinates.lat = lat;
        coordinates.lng = lng;
        zoom = 14;
    }
    // The map, centered at coordinates
    map = new google.maps.Map(document.getElementById('map'), {zoom: zoom, center: coordinates});
    //The marker, positioned at coordinates
    marker = new google.maps.Marker({position: coordinates, map: map});
}
function getLoaderHtml(){
    return '<div class="mx-auto my-5"><i class="fa fa-spinner fa-pulse fa-10x"></i></div>';
}