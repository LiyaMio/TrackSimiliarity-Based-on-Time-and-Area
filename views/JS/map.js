var drawingManager
var labelOptions
localStorage.setItem("evalue", "0.4")
localStorage.setItem("tvalue","20")
localStorage.setItem("stime","")
localStorage.setItem("etime","")
localStorage.setItem("method","")
var contrastData
function loadMapScenario(){
    $("#time-input").hide()
    $("#time-input2").hide()
    $(".modal").hide()
    var map = new BMapGL.Map('container'); // 创建Map实例
    map.centerAndZoom('上海市', 12); // 初始化地图,设置中心点坐标和地图级别
    map.enableScrollWheelZoom(true); // 开启鼠标滚轮缩放
    map.enableScrollWheelZoom(true);     //开启鼠标滚轮缩放
    var scaleCtrl = new BMapGL.ScaleControl();  // 添加比例尺控件
    map.addControl(scaleCtrl);
    var zoomCtrl = new BMapGL.ZoomControl();  // 添加比例尺控件
    map.addControl(zoomCtrl);
// 创建城市选择控件
    var cityControl = new BMapGL.CityListControl({
        // 控件的停靠位置（可选，默认左上角）
        anchor: BMAP_ANCHOR_TOP_LEFT,
        // 控件基于停靠位置的偏移量（可选）
        offset: new BMapGL.Size(10, 5)
    });
// 将控件添加到地图上
    map.addControl(cityControl);
// var marker = new BMapGL.Marker(point, {
//     enableDragging: true
// });
// map.addOverlay(marker);
    var styleOptions = {
        strokeColor: '#5E87DB',   // 边线颜色
        fillColor: '#5E87DB',     // 填充颜色。当参数为空时，圆形没有填充颜色
        strokeWeight: 2,          // 边线宽度，以像素为单位
        strokeOpacity: 1,         // 边线透明度，取值范围0-1
        fillOpacity: 0.2          // 填充透明度，取值范围0-1
    };
    labelOptions = {
        borderRadius: '2px',
        background: '#FFFBCC',
        border: '1px solid #E1E1E1',
        color: '#703A04',
        fontSize: '12px',
        letterSpacing: '0',
        padding: '5px'
    };
    drawingManager = new BMapGLLib.DrawingManager(map, {
        isOpen: false,        // 是否开启绘制模式
        enableCalculate: false, // 绘制是否进行测距测面
        enableSorption: true,   // 是否开启边界吸附功能
        sorptiondistance: 20,   // 边界吸附距离
        polylineOptions: styleOptions,   // 线的样式
    });
}
function draw(e,str) {
    var arr = document.getElementsByClassName('bmap-btn');
    for (var i = 0; i < arr.length; i++) {
        arr[i].style.backgroundPositionY = '0';
    }
    e.style.backgroundPositionY = '-52px';
    var drawingType = BMAP_DRAWING_POLYLINE;
    // 进行绘制
    if (drawingManager._isOpen && drawingManager.getDrawingMode() === drawingType) {
        drawingManager.close();
    } else {
        drawingManager.setDrawingMode(drawingType);
        drawingManager.open();
    }
    drawingManager.addEventListener("polylinecomplete", function (e) {
        $("#time-input").show()
        // $(".time-area").style.visibility = ($(".time-area").style.visibility == "visible") ? "hidden" : "visible";
        // var e1 = document.getElementById('time-input');
        // e1.style.visibility = (e1.style.visibility == "visible") ? "hidden" : "visible";

        path = e.getPath()

    });
};

var path = null

function timeConv(t) {
    var time = new Date(t);

    if(time.getMonth()+1>=10&&time.getDate()>=10)
        var strtim = time.getFullYear() + "-" + (time.getMonth()+1) + "-" + time.getDate() ;
    if(time.getMonth()+1>=10&&time.getDate()<10)
        var strtim = time.getFullYear() + "-" + (time.getMonth()+1) + "-0" + time.getDate() ;
    if(time.getMonth()+1<10&&time.getDate()<10)
        var strtim =time.getFullYear() + "-0" + (time.getMonth()+1) + "-0" + time.getDate() ;
    if(time.getMonth()+1<10&&time.getDate()>=10)
        var strtim = time.getFullYear() + "-0" + (time.getMonth()+1) + "-" + time.getDate() ;
    if(time.getHours()<10)
        strtim=strtim+" 0"+ time.getHours();
    else
        strtim=strtim+" "+ time.getHours();
    if(time.getMinutes()<10)
        strtim=strtim+ ":0" + time.getMinutes();
    else
        strtim=strtim+ ":" + time.getMinutes();
    if(time.getSeconds()<10)
        strtim=strtim+":0" + time.getSeconds();
    else
        strtim=strtim+":" + time.getSeconds();
    return strtim
}

function submitPath() {
    localStorage.setItem("stime",timeConv($("#stime")[0].value))
    localStorage.setItem("etime",timeConv($("#etime")[0].value))
    localStorage.setItem("method","path")
    $.ajax({
        type: "POST",
        url: '/map',
        headers: {
            'x-token': 'token'
        },
        data: JSON.stringify({
            stime: timeConv($("#stime")[0].value),
            etime: timeConv($("#etime")[0].value),
            path: path,
            距离阈值:localStorage.getItem("evalue"),
            时间阈值:localStorage.getItem("tvalue")

        }),
        dataType: 'JSON',
        contentType: 'application/json;charset=utf-8',
        async: false,
        success: function (data) {
            $(".navbar-nav").append("\n" +
                "                <li class=\"nav-item dropdown\">\n" +
                "                    <a class=\"nav-link dropdown-toggle\" href=\"#\" role=\"button\" data-bs-toggle=\"dropdown\" aria-expanded=\"false\">\n" +
                "                        阈值对比\n" +
                "                    </a>\n" +
                "                    <ul class=\"dropdown-menu\">\n" +
                "                        <li><a class=\"nav-link #thresholdModal\"  href=\"javascript:void(0)\" data-toggle=\"modal\" data-target=\"#thresholdModal\" onclick=\"contrast()\">阈值选择</a></li>\n" +
                "                        <li><a class=\"nav-link #contrastModal\"  href=\"javascript:void(0)\" data-toggle=\"modal\" data-target=\"#contrastModal\" >对比结果</a></li>\n" +
                "\n" +
                "                    </ul>\n" +
                "\n" +
                "                </li>")
            $("#time-input").hide()
            var map = new BMapGL.Map('container'); // 创建Map实例
            map.centerAndZoom('上海市', 12); // 初始化地图,设置中心点坐标和地图级别
            map.enableScrollWheelZoom(true); // 开启鼠标滚轮缩放
            map.enableScrollWheelZoom(true);     //开启鼠标滚轮缩放
            var scaleCtrl = new BMapGL.ScaleControl();  // 添加比例尺控件
            map.addControl(scaleCtrl);
            var zoomCtrl = new BMapGL.ZoomControl();  // 添加比例尺控件
            map.addControl(zoomCtrl);
// 创建城市选择控件
            var cityControl = new BMapGL.CityListControl({
                // 控件的停靠位置（可选，默认左上角）
                anchor: BMAP_ANCHOR_TOP_LEFT,
                // 控件基于停靠位置的偏移量（可选）
                offset: new BMapGL.Size(10, 5)
            });
// // 将控件添加到地图上
            map.addControl(cityControl);
            var nowponits=[]
            for (var i=0;i<path.length;i++){
                nowponits.push(new BMapGL.Point(path[i].lng, path[i].lat))
            }
            var nowline = new BMapGL.Polyline(nowponits, {strokeColor: "blue", strokeWeight: 2, strokeOpacity: 0.5});   //创建折线
            map.addOverlay(nowline);   //增加折线
            for(var i=0;i<data["len"];i++) {
                var len = data["track"]["N" + i].length
                var linePoints = []
                for (var j = 0; j < len; j++) {
                    linePoints.push(new BMapGL.Point(data["track"]["N" + i][j][0], data["track"]["N" + i][j][1]))
                }
               var polyline = new BMapGL.Polyline(linePoints, {strokeColor: "red", strokeWeight: 2, strokeOpacity: 0.5});   //创建折线
                map.addOverlay(polyline);   //增加折线
            }

        }
    })
}
function drawpath() {
    $.ajax({
        type: "POST",
        url: '/test',
        headers: {
            'x-token': 'token'
        },
        dataType: 'JSON',
        contentType: 'application/json;charset=utf-8',
        async: false,
        complete: function (data) {
            var res=data.track
            console.log(res)
            location.href=data.responseJSON.location
        }
    })
}
function regist() {
    if($("#registerConfirmInputPwd").val()!=$("#registerInputPwd").val()){
        alert("两次输入密码不一致")
        $("#registerConfirmInputPwd").append("<small id=\"emailHelp\" class=\"form-text text-muted\" style =\"color: darkred\"></small>")
    }
    else{
        $.ajax({
            type:'POST',
            url:"/register",
            dataType:'json',
            data:{
                name:$("#registerInputName").val(),
                pwd:$("#registerInputPwd").val(),
            },
            success:function (data) {
                if(data.msg=="success") {
                    alert("注册成功，请登录")
                    $("#register").attr("data-dismiss", "modal")
                }else{
                    alert(data.msg)
                }
            }
        })
    }
}
function login() {
    $.ajax({
        type:'POST',
        url:"/login",
        dataType:'json',
        data:{
            name:$("#exampleInputName").val(),
            pwd:$("#exampleInputPwd").val(),
        },
        success:function (data) {
            if(data.msg=="success") {
                alert("登录成功")
                $.ajax({
                    type:'POST',
                    url:"/cookie",
                    dataType:'json',
                    success:function (data) {
                        tracklist()
                    }
                })
            }else{
                alert(data.msg)
            }
        }
    })
}
var track
function AddTrack(e) {

    var arr = document.getElementsByClassName('bmap-btn');
    for (var i = 0; i < arr.length; i++) {
        arr[i].style.backgroundPositionY = '0';
    }
    e.style.backgroundPositionY = '-52px';
    var drawingType = BMAP_DRAWING_POLYLINE;
    // 进行绘制
    if (drawingManager._isOpen && drawingManager.getDrawingMode() === drawingType) {
        drawingManager.close();
    } else {
        drawingManager.setDrawingMode(drawingType);
        drawingManager.open();
    }
    var cnt=0
    // console.log(e.getPath().length)
    // console.log(cnt++)
    drawingManager.addEventListener("polylinecomplete", function (e) {
        $("#addtrackModal").show()

        track = e.getPath()
        console.log(track)


    });
    console.log(track)
}
function submitTrack() {
    $("#addtrackModal").hide()
    $.ajax({
        type: "POST",
        url: '/addtrack',
        headers: {
            'x-token': 'token'
        },
        data: JSON.stringify({
            tname:$("#addtrackInputName").val(),
            path: track
        }),
        dataType: 'JSON',
        contentType: 'application/json;charset=utf-8',
        async: false,
        success: function (data) {
            console.log(path)
        }
    })
}
function tracklist() {
    $(".tracklist").empty()
    $.ajax({
        type:'POST',
        url:"/postTrack",
        dataType:'json',
        success:function (data) {
            for(var i=0;i<data.length;i++){
                $(".tracklist").append("   <li class=\"list-group-item tracklistitem\" id=\""+data[i].Tid+"\"  onclick=\"active(this)\">"+data[i].Tname+"</li>")
            }
        }
    })

}
function active(e) {
    console.log(e.className.indexOf('active'))
    if (e.className.indexOf('active')==-1){
        var arr=document.getElementsByClassName("tracklistitem")
        for(var i=0;i<arr.length;i++){
            arr[i].className="list-group-item tracklistitem"
        }
        e.className+=" active"
    }else{
        var arr=document.getElementsByClassName("tracklistitem")
        for(var i=1;i<arr.length;i++){
            arr[i].className="list-group-item tracklistitem"
        }
        e.className="list-group-item tracklistitem"
    }
}
function chooseTrack() {
    var name
    var arr=document.getElementsByClassName("tracklistitem")
    for(var i=0;i<arr.length;i++){
        if (arr[i].className.indexOf('active')>0){
           name=arr[i].innerHTML
        }
    }
    localStorage.setItem("chooseName",name)
    console.log(name)
    $("#time-input2").show()
}
function submitTime() {
    $("#time-input").hide()
    $("#time-input2").hide()
    localStorage.setItem("stime", timeConv($("#sstime")[0].value))
    localStorage.setItem("etime",timeConv($("#eetime")[0].value))
    localStorage.setItem("method","common")
    $.ajax({
        type: "POST",
        url: '/choosetrack',
        headers: {
            'x-token': 'token'
        },
        data: JSON.stringify({
            stime: timeConv($("#sstime")[0].value),
            etime: timeConv($("#eetime")[0].value),
            name:localStorage.getItem("chooseName"),
            距离阈值:localStorage.getItem("evalue"),
            时间阈值:localStorage.getItem("tvalue")

        }),
        dataType: 'JSON',
        contentType: 'application/json;charset=utf-8',
        async: false,
        success: function (data) {
            $(".navbar-nav").append("\n" +
                "                <li class=\"nav-item dropdown\">\n" +
                "                    <a class=\"nav-link dropdown-toggle\" href=\"#\" role=\"button\" data-bs-toggle=\"dropdown\" aria-expanded=\"false\">\n" +
                "                        阈值对比\n" +
                "                    </a>\n" +
                "                    <ul class=\"dropdown-menu\">\n" +
                "                        <li><a class=\"nav-link #thresholdModal\"  href=\"javascript:void(0)\" data-toggle=\"modal\" data-target=\"#thresholdModal\" onclick=\"contrast()\">阈值选择</a></li>\n" +
                "                        <li><a class=\"nav-link #contrastModal\"  href=\"javascript:void(0)\" data-toggle=\"modal\" data-target=\"#contrastModal\" >对比结果</a></li>\n" +
                "\n" +
                "                    </ul>\n" +
                "\n" +
                "                </li>")
            var map = new BMapGL.Map('container'); // 创建Map实例
            map.centerAndZoom('上海市', 12); // 初始化地图,设置中心点坐标和地图级别
            map.enableScrollWheelZoom(true); // 开启鼠标滚轮缩放
            map.enableScrollWheelZoom(true);     //开启鼠标滚轮缩放
            var scaleCtrl = new BMapGL.ScaleControl();  // 添加比例尺控件
            map.addControl(scaleCtrl);
            var zoomCtrl = new BMapGL.ZoomControl();  // 添加比例尺控件
            map.addControl(zoomCtrl);
// 创建城市选择控件
            var cityControl = new BMapGL.CityListControl({
                // 控件的停靠位置（可选，默认左上角）
                anchor: BMAP_ANCHOR_TOP_LEFT,
                // 控件基于停靠位置的偏移量（可选）
                offset: new BMapGL.Size(10, 5)
            });
// 将控件添加到地图上
            map.addControl(cityControl);
            var nowponits=[]
            for (var i=0;i<data["orgin"].length;i++){
                nowponits.push(new BMapGL.Point(data["orgin"][i]["lng"], data["orgin"][i]["lat"]))
            }
            var nowline = new BMapGL.Polyline(nowponits, {strokeColor: "blue", strokeWeight: 2, strokeOpacity: 0.5});   //创建折线
            map.addOverlay(nowline);   //增加折线
            for(var i=0;i<data["len"];i++) {
                var len = data["track"]["N" + i].length
                var linePoints = []
                for (var j = 0; j < len; j++) {
                    linePoints.push(new BMapGL.Point(data["track"]["N" + i][j][0], data["track"]["N" + i][j][1]))
                }
                var polyline = new BMapGL.Polyline(linePoints, {strokeColor: "red", strokeWeight: 2, strokeOpacity: 0.5});   //创建折线
                map.addOverlay(polyline);   //增加折线
            }

        }
    })
}
function getE(e) {
    if (e.text == "默认阈值") {
        localStorage.setItem("evalue", "0.4")
    } else {
        localStorage.setItem("evalue", e.text)
    }
    $("#thresholdD").attr("placeholder","当前阈值为："+localStorage.getItem("evalue"))
}
function getT(e) {
    if(e.text=="默认阈值"){
        localStorage.setItem("tvalue","20")
    }else{
        localStorage.setItem("tvalue",e.text)
    }
    $("#thresholdT").attr("placeholder","当前阈值为："+localStorage.getItem("tvalue"))
}
function contrastResult() {
    console.log(contrastData.length)
    var map = new BMapGL.Map('contrastmap'); // 创建Map实例
    map.centerAndZoom('上海市', 12); // 初始化地图,设置中心点坐标和地图级别
    map.enableScrollWheelZoom(true); // 开启鼠标滚轮缩放
    map.enableScrollWheelZoom(true);     //开启鼠标滚轮缩放
    var scaleCtrl = new BMapGL.ScaleControl();  // 添加比例尺控件
    map.addControl(scaleCtrl);
    var zoomCtrl = new BMapGL.ZoomControl();  // 添加比例尺控件
    map.addControl(zoomCtrl);
// 创建城市选择控件
    var cityControl = new BMapGL.CityListControl({
        // 控件的停靠位置（可选，默认左上角）
        anchor: BMAP_ANCHOR_TOP_LEFT,
        // 控件基于停靠位置的偏移量（可选）
        offset: new BMapGL.Size(10, 5)
    });
// 将控件添加到地图上
    map.addControl(cityControl);
// var marker = new BMapGL.Marker(point, {
//     enableDragging: true
// });
// map.addOverlay(marker);
    var styleOptions = {
        strokeColor: '#5E87DB',   // 边线颜色
        fillColor: '#5E87DB',     // 填充颜色。当参数为空时，圆形没有填充颜色
        strokeWeight: 2,          // 边线宽度，以像素为单位
        strokeOpacity: 1,         // 边线透明度，取值范围0-1
        fillOpacity: 0.2          // 填充透明度，取值范围0-1
    };
    labelOptions = {
        borderRadius: '2px',
        background: '#FFFBCC',
        border: '1px solid #E1E1E1',
        color: '#703A04',
        fontSize: '12px',
        letterSpacing: '0',
        padding: '5px'
    };
    drawingManager = new BMapGLLib.DrawingManager(map, {
        isOpen: false,        // 是否开启绘制模式
        enableCalculate: false, // 绘制是否进行测距测面
        enableSorption: true,   // 是否开启边界吸附功能
        sorptiondistance: 20,   // 边界吸附距离
        polylineOptions: styleOptions,   // 线的样式
    });
}
function onlyNumber(obj){
    //得到第一个字符是否为负号
    var t = obj.value.charAt(0);
    //先把非数字的都替换掉，除了数字和.
    obj.value = obj.value.replace(/[^\d\.]/g,'');
    //必须保证第一个为数字而不是.
    obj.value = obj.value.replace(/^\./g,'');
    //保证只有出现一个.而没有多个.
    obj.value = obj.value.replace(/\.{2,}/g,'.');
    //保证.只出现一次，而不能出现两次以上
    obj.value = obj.value.replace('.','$#$').replace(/\./g,'').replace('$#$','.');
    //如果第一位是负号，则允许添加
    if(t == '-'){
        obj.value = '-'+obj.value;
    }
}
function contrast() {
    if(localStorage.getItem("method")=="path") {
        $.ajax({
            type: "POST",
            url: '/contrastpath',
            headers: {
                'x-token': 'token'
            },
            data: JSON.stringify({
                stime: localStorage.getItem("stime"),
                etime: localStorage.getItem("etime"),
                path: path,
                距离阈值: $("#thresholdD").val(),
                时间阈值: $("#thresholdT").val()

            }),
            dataType: 'JSON',
            contentType: 'application/json;charset=utf-8',
            async: false,
            success: function (data) {
                var map = new BMapGL.Map('contrastmap'); // 创建Map实例
                map.centerAndZoom('上海市', 12); // 初始化地图,设置中心点坐标和地图级别
                map.enableScrollWheelZoom(true); // 开启鼠标滚轮缩放
                map.enableScrollWheelZoom(true);     //开启鼠标滚轮缩放
                var scaleCtrl = new BMapGL.ScaleControl();  // 添加比例尺控件
                map.addControl(scaleCtrl);
                var zoomCtrl = new BMapGL.ZoomControl();  // 添加比例尺控件
                map.addControl(zoomCtrl);
// 创建城市选择控件
                var cityControl = new BMapGL.CityListControl({
                    // 控件的停靠位置（可选，默认左上角）
                    anchor: BMAP_ANCHOR_TOP_LEFT,
                    // 控件基于停靠位置的偏移量（可选）
                    offset: new BMapGL.Size(10, 5)
                });
// // 将控件添加到地图上
                map.addControl(cityControl);
                var nowponits=[]
                for (var i=0;i<path.length;i++){
                    nowponits.push(new BMapGL.Point(path[i].lng, path[i].lat))
                }
                var nowline = new BMapGL.Polyline(nowponits, {strokeColor: "blue", strokeWeight: 2, strokeOpacity: 0.5});   //创建折线
                map.addOverlay(nowline);   //增加折线
                for(var i=0;i<data["len"];i++) {
                    var len = data["track"]["N" + i].length
                    var linePoints = []
                    for (var j = 0; j < len; j++) {
                        linePoints.push(new BMapGL.Point(data["track"]["N" + i][j][0], data["track"]["N" + i][j][1]))
                    }
                    var polyline = new BMapGL.Polyline(linePoints, {strokeColor: "red", strokeWeight: 2, strokeOpacity: 0.5});   //创建折线
                    map.addOverlay(polyline);   //增加折线

                }
            }
        })
    }
    else{

        $.ajax({
            type: "POST",
            url: '/contrastcommon',
            headers: {
                'x-token': 'token'
            },
            data: JSON.stringify({
                stime: localStorage.getItem("stime"),
                etime: localStorage.getItem("etime"),
                name:localStorage.getItem("chooseName"),
                距离阈值: $("#thresholdD").val(),
                时间阈值: $("#thresholdT").val()

            }),
            dataType: 'JSON',
            contentType: 'application/json;charset=utf-8',
            async: false,
            success: function (data) {
                contrastData=data

                var map = new BMapGL.Map('contrastmap'); // 创建Map实例
                map.centerAndZoom('上海市', 12); // 初始化地图,设置中心点坐标和地图级别
                map.enableScrollWheelZoom(true); // 开启鼠标滚轮缩放
                map.enableScrollWheelZoom(true);     //开启鼠标滚轮缩放
                var scaleCtrl = new BMapGL.ScaleControl();  // 添加比例尺控件
                map.addControl(scaleCtrl);
                var zoomCtrl = new BMapGL.ZoomControl();  // 添加比例尺控件
                map.addControl(zoomCtrl);
// 创建城市选择控件
                var cityControl = new BMapGL.CityListControl({
                    // 控件的停靠位置（可选，默认左上角）
                    anchor: BMAP_ANCHOR_TOP_LEFT,
                    // 控件基于停靠位置的偏移量（可选）
                    offset: new BMapGL.Size(10, 5)
                });
// 将控件添加到地图上
                map.addControl(cityControl);
                var nowponits = []
                for (var i = 0; i < data["orgin"].length; i++) {
                    nowponits.push(new BMapGL.Point(data["orgin"][i]["lng"], data["orgin"][i]["lat"]))
                }
                var nowline = new BMapGL.Polyline(nowponits, {
                    strokeColor: "blue",
                    strokeWeight: 2,
                    strokeOpacity: 0.5
                });   //创建折线
                map.addOverlay(nowline);   //增加折线
                for (var i = 0; i < data["len"]; i++) {
                    var len = data["track"]["N" + i].length
                    var linePoints = []
                    for (var j = 0; j < len; j++) {
                        linePoints.push(new BMapGL.Point(data["track"]["N" + i][j][0], data["track"]["N" + i][j][1]))
                    }
                    var polyline = new BMapGL.Polyline(linePoints, {
                        strokeColor: "red",
                        strokeWeight: 2,
                        strokeOpacity: 0.5
                    });   //创建折线
                    map.addOverlay(polyline);   //增加折线
                }
            }
        })
    }
}
function deleteTrack() {
    var name
    var arr=document.getElementsByClassName("tracklistitem")
    for(var i=0;i<arr.length;i++){
        if (arr[i].className.indexOf('active')>0){
            name=arr[i].innerHTML
        }
    }
    $.ajax({
        type: "POST",
        url: '/deletetrack',
        headers: {
            'x-token': 'token'
        },
        data: JSON.stringify({
            name:name
        }),
        success:function (data) {
            alert(data.msg)
        }
    })
}
function search() {
    var map = new BMapGL.Map("container");
    map.centerAndZoom(new BMapGL.Point(116.404, 39.915), 11);
    var local = new BMapGL.LocalSearch(map, {
        renderOptions:{map: map}
    });
    local.search("景点");
}