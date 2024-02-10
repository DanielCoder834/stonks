import * as data from "./svgs/Animation - 1704740731596.json";

let pageX = 0 
let pageY = 0
let circle = document.querySelector('.circle');
circle.style.position = 'fixed';  

let cat = document.getElementById('cat');
console.log(typeof cat);
let mydata = JSON.parse(data);
console.log(mydata);
// let catHeight = cat.clientHeight;
// let catWidth = cat.clientWidth;

onmousemove = function(e) {
    // console.log("cat height and width: ", catHeight, catWidth);
    // console.log("mouse location:", e.clientX, e.clientY)
    pageX = e.clientX - (circle.offsetWidth * 0.75)
    pageY = e.clientY - (circle.offsetHeight * 0.75)
    // console.log("page locations: ", pageX, pageY)
    anime({
        targets: '.circle',
        translateX: pageX,
        translateY: pageY,
        // delay: anime.stagger(200, {start: 1000}),
        // background: "#0000FF",
        // direction: "reverse",
        // duration: 1000,
    })
}


