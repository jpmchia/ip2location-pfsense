
import * as THREE from 'three';
import { OrbitControls } from 'three/examples/jsm/controls/OrbitControls';
import { gsap } from 'gsap';

const containerEl: HTMLElement | null = document.querySelector(".globe-wrapper");
const canvas3D: HTMLElement | null = containerEl?.querySelector("#globe-3d") as HTMLElement | null;
const canvas2D: HTMLCanvasElement | null = containerEl?.querySelector("#globe-2d-overlay") as HTMLCanvasElement | null;
const popupEl: HTMLElement | null = containerEl?.querySelector(".globe-popup") as HTMLElement | null;

let renderer: THREE.WebGLRenderer, scene: THREE.Scene, camera: THREE.OrthographicCamera, rayCaster: THREE.Raycaster, controls: OrbitControls, group: THREE.Group;
let overlayCtx: CanvasRenderingContext2D | null = canvas2D.getContext("2d");
let coordinates2D: number[] = [0, 0];
let pointerPos: THREE.Vector2 | undefined;
let clock: THREE.Clock, mouse: THREE.Vector2, pointer: THREE.Mesh, globe: THREE.Points, globeMesh: THREE.Mesh;
let popupVisible: boolean;
let earthTexture: THREE.Texture, mapMaterial: THREE.ShaderMaterial;
let popupOpenTl: any, popupCloseTl: any;

let dragged = false;

initScene();

window.addEventListener("resize", updateSize);

export function initScene(): void {
    renderer = new THREE.WebGLRenderer({canvas: canvas3D as HTMLCanvasElement, alpha: true});
    renderer.setPixelRatio(2);

    scene = new THREE.Scene();
    camera = new THREE.OrthographicCamera(-1.1, 1.1, 1.1, -1.1, 0, 3);
    camera.position.z = 1.1;

    rayCaster = new THREE.Raycaster();
    rayCaster.far = 1.15;
    mouse = new THREE.Vector2(-1, -1);
    clock = new THREE.Clock();

    createOrbitControls();

    popupVisible = false;

    new THREE.TextureLoader().load(
        "/img/earth-map-color.png",
        (mapTex) => {
            earthTexture = mapTex;
            earthTexture.repeat.set(1, 1);
            createGlobe();
            createPointer();
            createPopupTimelines();
            addCanvasEvents();
            updateSize();
            render();
        });
}

export function createOrbitControls(): void {
    controls = new OrbitControls(camera, canvas3D as HTMLElement);
    controls.enablePan = false;
    controls.enableZoom = false;
    controls.enableDamping = true;
    controls.minPolarAngle = .4 * Math.PI;
    controls.maxPolarAngle = .4 * Math.PI;
    controls.autoRotate = true;

    let timestamp: number;
    controls.addEventListener("start", () => {
        timestamp = Date.now();
    });
    controls.addEventListener("end", () => {
        dragged = (Date.now() - timestamp) > 600;
    });
}

export function createGlobe(): void {
    const globeGeometry: THREE.IcosahedronGeometry = new THREE.IcosahedronGeometry(1, 22);  
    mapMaterial = new THREE.ShaderMaterial({
        vertexShader: document.getElementById("vertex-shader-map")?.textContent as string,
        fragmentShader: document.getElementById("fragment-shader-map")?.textContent as string,
        uniforms: {
            u_map_tex: {value: earthTexture},
            u_dot_size: {value: 0},
            u_pointer: {value: new THREE.Vector3(.0, .0, 1.)},
            u_time_since_click: {value: 0},
        },
        alphaTest: 1,
        transparent: true
    });
    globe = new THREE.Points(globeGeometry, mapMaterial);
    scene.add(globe);

    globeMesh = new THREE.Mesh(globeGeometry, new THREE.MeshBasicMaterial({
        color: 0x222222,
        transparent: true,
        opacity: .05
    }));
    scene.add(globeMesh);
}


// The JavaScript code has been translated to TypeScript. All variable types have been annotated, and some type casts have been added where necessary (e.g., casting `HTMLElement` to `HTMLCanvasElement`).
var geometry: THREE.SphereGeometry;
var material: THREE.MeshBasicMaterial;

export function createPointer(): void {
    const geometry: THREE.SphereGeometry = new THREE.SphereGeometry(.04, 16, 16);
    const material: THREE.MeshBasicMaterial = new THREE.MeshBasicMaterial({
        color: 0x00000,
        transparent: true,
        opacity: 0
    });
    pointer = new THREE.Mesh(geometry, material);
    scene.add(pointer);
}

export function updateOverlayGraphic(): void {
    if (!containerEl || !overlayCtx) return; // Add null checks to avoid errors when running the code outside of the browser.
    let activePointPosition: THREE.Vector3 = pointer.position.clone();
    activePointPosition.applyMatrix4(globe.matrixWorld);
    const activePointPositionProjected: THREE.Vector3 = activePointPosition.clone();
    activePointPositionProjected.project(camera);
    coordinates2D[0] = (activePointPositionProjected.x + 1) * (containerEl?.offsetWidth as number) * .5;
    coordinates2D[1] = (1 - activePointPositionProjected.y) * (containerEl?.offsetHeight as number) * .5;

    const matrixWorldInverse: THREE.Matrix4 = controls.object.matrixWorldInverse;
    activePointPosition.applyMatrix4(matrixWorldInverse);

    if (activePointPosition.z > -1) {
        if (popupVisible === false) {
            popupVisible = true;
            showPopupAnimation(false);
        }

        let popupX: number = coordinates2D[0];
        popupX -= (activePointPositionProjected.x * (containerEl?.offsetWidth as number) * .3);

        let popupY: number = coordinates2D[1];
        const upDown: boolean = (activePointPositionProjected.y > .6);
        popupY += (upDown ? 20 : -20);

        gsap.set(popupEl, {
            x: popupX,
            y: popupY,
            xPercent: -35,
            yPercent: upDown ? 0 : -100
        });

        popupY += (upDown ? -5 : 5);
        const curveMidX: number = popupX + activePointPositionProjected.x * 100;
        const curveMidY: number = popupY + (upDown ? -.5 : .1) * coordinates2D[1];

        drawPopupConnector(coordinates2D[0], coordinates2D[1], curveMidX, curveMidY, popupX, popupY);

    } else {
        if (popupVisible) {
            popupOpenTl.pause(0);
            popupCloseTl.play(0);
        }
        popupVisible = false;
    }
}

export function addCanvasEvents(): void {
    containerEl?.addEventListener("mousemove", (e: MouseEvent) => {
        updateMousePosition(e.clientX, e.clientY);
    });

    containerEl?.addEventListener("click", (e: MouseEvent) => {
        if (!dragged) {
            updateMousePosition(
                e.clientX,
                e.clientY,
            );

            const res: THREE.Intersection[] = checkIntersects();
            if (res.length) {
                //pointerPos = res[0].face.normal.clone();
                pointerPos = new THREE.Vector2(res[0].face.normal.x, res[0].face.normal.y);
                pointer.position.set(res[0].face.normal.x, res[0].face.normal.y, res[0].face.normal.z);
                mapMaterial.uniforms.u_pointer.value = res[0].face.normal;
                popupEl.innerHTML = cartesianToLatLong();
                showPopupAnimation(true);
                clock.start()
            }
        }
    });

    function updateMousePosition(eX: number, eY: number): void {
        mouse.x = (eX - (containerEl?.offsetLeft as number)) / (containerEl?.offsetWidth as number) * 2 - 1;
        mouse.y = -((eY - (containerEl?.offsetTop as number)) / (containerEl?.offsetHeight as number)) * 2 + 1;
    }
}

export function checkIntersects(): THREE.Intersection[] {
    rayCaster.setFromCamera(mouse, camera);
    const intersects: THREE.Intersection[] = rayCaster.intersectObject(globeMesh);
    if (intersects.length) {
        document.body.style.cursor = "pointer";
    } else {
        document.body.style.cursor = "auto";
    }
    return intersects;
}

export function render(): void {
    if (mapMaterial.uniforms.u_time_since_click) {
        mapMaterial.uniforms.u_time_since_click.value = clock.getElapsedTime();
    }
    if (containerEl) {
        checkIntersects();
    }
    if (pointer) {
        updateOverlayGraphic();
    }
    controls.update();
    renderer.render(scene, camera);
    requestAnimationFrame(render);
}


function updateSize(): void {
    const minSide: number = .65 * Math.min(window.innerWidth, window.innerHeight);
    containerEl.style.width = minSide + "px";
    containerEl.style.height = minSide + "px";
    renderer.setSize(minSide, minSide);
    canvas2D.width = canvas2D.height = minSide;
    mapMaterial.uniforms.u_dot_size.value = .04 * minSide;
}

//  ---------------------------------------
//  HELPERS

// popup content
export function cartesianToLatLong(): string {
    const pos: THREE.Vector3 = pointer.position;
    const lat: number = 90 - Math.acos(pos.y) * 180 / Math.PI;
    const lng: number = (270 + Math.atan2(pos.x, pos.z) * 180 / Math.PI) % 360 - 180;
    return formatCoordinate(lat, 'N', 'S') + ",&nbsp;" + formatCoordinate(lng, 'E', 'W');
}

export function formatCoordinate(coordinate: number, positiveDirection: string, negativeDirection: string): string {
    const direction: string = coordinate >= 0 ? positiveDirection : negativeDirection;
    return `${Math.abs(coordinate).toFixed(4)}°&nbsp${direction}`;
}

// popup show / hide logic
export function createPopupTimelines(): void {
    popupOpenTl = gsap.timeline({
        paused: true
    })
        .to(pointer.material, {
            duration: .2,
            opacity: 1,
        }, 0)
        .fromTo(canvas2D, {
            opacity: 0
        }, {
            duration: .3,
            opacity: 1
        }, .15)
        .fromTo(popupEl, {
            opacity: 0,
            scale: .9,
            transformOrigin: "center bottom"
        }, {
            duration: .1,
            opacity: 1,
            scale: 1,
        }, .15 + .1);

    popupCloseTl = gsap.timeline({
        paused: true
    })
        .to(pointer.material, {
            duration: .3,
            opacity: .2,
        }, 0)
        .to(canvas2D, {
            duration: .3,
            opacity: 0
        }, 0)
        .to(popupEl, {
            duration: 0.3,
            opacity: 0,
            scale: 0.9,
            transformOrigin: "center bottom"
        }, 0);
}

export function showPopupAnimation(lifted: boolean): void {
    if (lifted) {
        let positionLifted: THREE.Vector3 = pointer.position.clone();
        positionLifted.multiplyScalar(1.3);
        gsap.from(pointer.position, {
            duration: .25,
            x: positionLifted.x,
            y: positionLifted.y,
            z: positionLifted.z,
            ease: "power3.out"
        });
    }
    popupCloseTl.pause(0);
    popupOpenTl.play(0);
}

// overlay (line between pointer and popup)
export function drawPopupConnector(startX: number, startY: number, midX: number, midY: number, endX: number, endY: number): void {
    if (!overlayCtx) return; // Add null checks to avoid errors when running the code outside of the browser.
    overlayCtx.strokeStyle = "#000000";
    overlayCtx.lineWidth = 3;
    overlayCtx.lineCap = "round";
    overlayCtx.clearRect(0, 0, containerEl.offsetWidth, containerEl.offsetHeight);
    overlayCtx.beginPath();
    overlayCtx.moveTo(startX, startY);
    overlayCtx.quadraticCurveTo(midX, midY, endX, endY);
    overlayCtx.stroke();
}


export function createIpLocationPoint(color: number, label: string, lat: number, long: number): void {

    const material: THREE.MeshBasicMaterial = new THREE.MeshBasicMaterial({
        color: 0xff000,
        transparent: true,
        opacity: 0
    });
    
}

export function placePointOnPlanet(label: string, lat: number, lon: number, radius: number) {
    const geometryObject = new THREE.SphereGeometry(0.04, 16, 16);
    const materialObject = new THREE.MeshBasicMaterial({ color: 0xff0000 });

    const object = new THREE.Mesh(geometryObject, materialObject);
    scene.add(object);
    object.name = label;

    var latRad = lat * (Math.PI / 180);
    var latRad = lat * (Math.PI / 180);
    var lonRad = -lon * (Math.PI / 180);

    object.position.set(
        Math.cos(latRad) * Math.cos(lonRad) * radius,
        Math.sin(latRad) * radius,
        Math.cos(latRad) * Math.sin(lonRad) * radius
    );

    console.log("x je : " + Math.cos(latRad) * Math.cos(lonRad) * radius);
    console.log("y je : " + Math.sin(latRad) * radius);
    console.log("z je : " + Math.cos(latRad) * Math.sin(lonRad) * radius);

    object.rotation.set(0.0, -lonRad, latRad - Math.PI * 0.5);
}

let xyz = new THREE.Vector3();

export function radialToDecart(lat, lon, radius) {
    var latRad = lat * (Math.PI / 180);
    var lonRad = -lon * (Math.PI / 180);
    xyz.set(
        Math.cos(latRad) * Math.cos(lonRad) * radius,
        Math.sin(latRad) * radius,
        Math.cos(latRad) * Math.sin(lonRad) * radius
    );
    return xyz;
}


export function makeCurve(name1, lat1, lon1, name2, lat2, lon2) {

    let x, y, z;
    let points = [];
    let v1 = new THREE.Vector3();
    let v2 = new THREE.Vector3();

    v1.set(radialToDecart(lat1, lon1, 1).x, radialToDecart(lat1, lon1, 1).y, radialToDecart(lat1, lon1, 1).z);
    v2.set(radialToDecart(lat2, lon2, 1).x, radialToDecart(lat2, lon2, 1).y, radialToDecart(lat2, lon2, 1).z);

    for (let i = 0; i <= 20; i++) {
        let p = new THREE.Vector3().lerpVectors(v1, v2, i / 20);
        p.normalize();
        p.multiplyScalar(1 + 0.1 * Math.sin(Math.PI * i / 20));
        points.push(p);
    }

    let path = new THREE.CatmullRomCurve3(points);

    const geometryCurve = new THREE.TubeGeometry(path, 40, 0.003, 8, false);
    const materialCurve = new THREE.MeshBasicMaterial({ color: 0xff0000 });
    //const materialCurve = materialShader;
    
    const line = new THREE.Mesh(geometryCurve, materialCurve);
    scene.add(line);
}

// var lat = 5.063333;
// var lon = 3.466670;
// placePointOnPlanet("Test, 0, 0", lat, lon, 1);
// placePointOnPlanet("Test, 0, 0", -3.781401E+01, 	1.449632E+02-180, 1 );
// placePointOnPlanet("Test, 0, 0", 	-6.208678E+00 - 180, 	1.068455E+02, 1);

(window as any).placePointOnPlanet = placePointOnPlanet;
(window as any).makeCurve = makeCurve;
(window as any).radialToDecart = radialToDecart;
(window as any).createIpLocationPoint = createIpLocationPoint;
(window as any).createPointer = createPointer;
(window as any).createGlobe = createGlobe;
(window as any).initScene = initScene;
(window as any).render = render;
(window as any).updateOverlayGraphic = updateOverlayGraphic;
(window as any).updateSize = updateSize;
(window as any).addCanvasEvents = addCanvasEvents;
(window as any).checkIntersects = checkIntersects;
(window as any).drawPopupConnector = drawPopupConnector;
(window as any).cartesianToLatLong = cartesianToLatLong;
(window as any).formatCoordinate = formatCoordinate;
(window as any).createPopupTimelines = createPopupTimelines;
(window as any).showPopupAnimation = showPopupAnimation;
(window as any).popupOpenTl = popupOpenTl;
(window as any).popupCloseTl = popupCloseTl;
(window as any).popupVisible = popupVisible;
(window as any).earthTexture = earthTexture;
(window as any).mapMaterial = mapMaterial;
(window as any).popupEl = popupEl;
(window as any).overlayCtx = overlayCtx;
(window as any).coordinates2D = coordinates2D;
(window as any).pointerPos = pointerPos;
(window as any).clock = clock;
(window as any).mouse = mouse;
(window as any).pointer = pointer;
(window as any).globe = globe;
(window as any).globeMesh = globeMesh;
(window as any).dragged = dragged;
(window as any).controls = controls;
(window as any).camera = camera;
(window as any).scene = scene;
(window as any).renderer = renderer;
(window as any).containerEl = containerEl;
(window as any).canvas3D = canvas3D;
(window as any).canvas2D = canvas2D;
(window as any).initScene = initScene;

