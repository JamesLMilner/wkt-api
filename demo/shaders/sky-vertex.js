
var Lacuna = (function () {

  // Instance stores a reference to the Singleton
  var instance;

  function startGame() {

    // Singleton

	var camera, scene, renderer;
	var lightIntensity = 0.15;
	var fog = 700; // Higher = less fog
    var LACUNAAPI = "localhost:8080/lacuna";

	init();
	animate();

	var prevTime = performance.now();
	var velocity = new THREE.Vector3();


	function init() {

		eventHandlers();
		scene = new THREE.Scene();
		scene.fog = new THREE.Fog( 0xffffff, 0, fog + 1000 );

		// Sky
		var pwd = window.location.href.substring(0, window.location.href.indexOf('/'));
		var sky = new THREE.SphereGeometry(8000, 32, 32); // radius, widthSegments, heightSegments
		var uniforms = {
		  texture: { type: 't', value: THREE.ImageUtils.loadTexture(pwd + 'imgs/sky.jpg') }
		};

		var skyMaterial = new THREE.ShaderMaterial( {
			uniforms:       uniforms,
			vertexShader:   document.getElementById('sky-vertex').textContent,
			fragmentShader: document.getElementById('sky-fragment').textContent
		});

		skyBox = new THREE.Mesh(sky, skyMaterial);
		skyBox.scale.set(-1, 1, 1);
		skyBox.eulerOrder = 'XZY';
		skyBox.renderDepth = 1000.0;
		scene.add(skyBox);

		camera = new THREE.PerspectiveCamera( 80, window.innerWidth / window.innerHeight, 1, 9000 );
		controls = new THREE.PointerLockControls( camera );
		scene.add( controls.getObject() );

		renderer = new THREE.WebGLRenderer({ antialias: true }); //new THREE.WebGLRenderer();
		renderer.setClearColor( 0xffffff );
		renderer.setPixelRatio( window.devicePixelRatio );
		renderer.setSize( window.innerWidth, window.innerHeight );
		document.body.appendChild( renderer.domElement );

	}

	function animate() {

		requestAnimationFrame( animate );
		renderer.render( scene, camera );

	}

	function eventHandlers() {

	}

    function getAllGeometries(table) {


    }


    function getGeometryById(id, table) {
        $.get()

    }



	function onWindowResize() {

		camera.aspect = window.innerWidth / window.innerHeight;
		camera.updateProjectionMatrix();

		renderer.setSize( window.innerWidth, window.innerHeight );

	}

    return {
		// Public methods and variables
		setFog: function (setFog) {
			fog = setFog;
		},
		setJumpFactor: function (setJumpFactor) {
			jumpFactor = setJumpFactor;
		}

    };

  }

  return {

    // Get the Singleton instance if one exists
    // or create one if it doesn't
    getInstance: function () {

      if ( !instance ) {
        instance = startGame();
      }

      return instance;
    }

  };

})();

Lacuna = Lacuna.getInstance();
