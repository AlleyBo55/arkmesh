"use client";

import { useEffect, useRef, useState } from "react";

const chapters = [
  {
    depth: "000 M",
    eyebrow: "00 / Dependency",
    title: "What survives when access disappears?",
    body: "A useful model that people are legally allowed to preserve can vanish from practical reach when its host, runtime, documentation, or trust evidence disappears. Downloading weights alone does not preserve the capability.",
  },
  {
    depth: "0240 M",
    eyebrow: "01 / Complete object",
    title: "Preserve everything needed to recover.",
    body: "ArkMesh packages model files, tokenizer data, configuration, runtime requirements, knowledge, licenses, lineage, and recovery instructions into one verifiable capsule.",
  },
  {
    depth: "1800 M",
    eyebrow: "02 / Current proof",
    title: "Exact bytes can return from pieces.",
    body: "The current release reconstructs an exact 8 MiB signed object from four of six shards after the publisher, complete object, and two shards are deleted.",
  },
  {
    depth: "4200 M",
    eyebrow: "03 / Trust boundary",
    title: "Recovery is accepted only after verification.",
    body: "Digests, Merkle commitments, signatures, lineage rules, revocation policy, and retained checkpoints decide whether recovered bytes may become an accepted capsule.",
  },
  {
    depth: "6000 M",
    eyebrow: "04 / Open experiment",
    title: "The peer network must still be built.",
    body: "ArkMesh does not yet discover peers or run models. The next experiment is consent-based encrypted recovery across independent devices without a central tracker.",
  },
] as const;

const chapterStarts = [0, 0.2, 0.41, 0.63, 0.82];

function getChapter(progress: number) {
  for (let index = chapterStarts.length - 1; index >= 0; index -= 1) {
    if (progress >= chapterStarts[index]) {
      return index;
    }
  }
  return 0;
}

export function OceanJourney() {
  const sectionRef = useRef<HTMLElement>(null);
  const canvasHostRef = useRef<HTMLDivElement>(null);
  const progressRef = useRef(0);
  const depthRef = useRef<HTMLSpanElement>(null);
  const [activeChapter, setActiveChapter] = useState(0);

  useEffect(() => {
    const section = sectionRef.current;
    if (!section) {
      return;
    }

    let frame = 0;
    const updateProgress = () => {
      frame = 0;
      const bounds = section.getBoundingClientRect();
      const travel = Math.max(section.offsetHeight - window.innerHeight, 1);
      const progress = Math.min(1, Math.max(0, -bounds.top / travel));
      progressRef.current = progress;
      section.style.setProperty("--dive-progress", progress.toFixed(4));

      if (depthRef.current) {
        const depth = Math.round(progress * 6000);
        depthRef.current.textContent = `${depth.toString().padStart(4, "0")} M`;
      }

      const nextChapter = getChapter(progress);
      setActiveChapter((current) => current === nextChapter ? current : nextChapter);
    };

    const requestUpdate = () => {
      if (!frame) {
        frame = window.requestAnimationFrame(updateProgress);
      }
    };

    updateProgress();
    window.addEventListener("scroll", requestUpdate, { passive: true });
    window.addEventListener("resize", requestUpdate);

    return () => {
      window.removeEventListener("scroll", requestUpdate);
      window.removeEventListener("resize", requestUpdate);
      if (frame) {
        window.cancelAnimationFrame(frame);
      }
    };
  }, []);

  useEffect(() => {
    const host = canvasHostRef.current;
    const section = sectionRef.current;
    const motionQuery = window.matchMedia("(prefers-reduced-motion: reduce)");

    if (!host || !section || motionQuery.matches) {
      return;
    }

    let cancelled = false;
    let disposeScene: (() => void) | undefined;

    async function mountScene() {
      const THREE = await import("three");
      if (cancelled || !host || !section) {
        return;
      }

      const scene = new THREE.Scene();
      const surfaceColor = new THREE.Color(0x162d48);
      const abyssColor = new THREE.Color(0x03050b);
      const currentColor = new THREE.Color();
      scene.background = currentColor.copy(surfaceColor);
      const fog = new THREE.FogExp2(surfaceColor, 0.038);
      scene.fog = fog;

      const camera = new THREE.PerspectiveCamera(46, 1, 0.1, 120);
      const renderer = new THREE.WebGLRenderer({
        alpha: false,
        antialias: false,
        powerPreference: "high-performance",
      });
      renderer.outputColorSpace = THREE.SRGBColorSpace;
      renderer.toneMapping = THREE.ACESFilmicToneMapping;
      renderer.toneMappingExposure = 1.12;
      renderer.setPixelRatio(Math.min(window.devicePixelRatio, window.innerWidth < 768 ? 1 : 1.5));
      renderer.domElement.setAttribute("aria-hidden", "true");
      host.appendChild(renderer.domElement);

      const world = new THREE.Group();
      scene.add(world);

      const ambient = new THREE.HemisphereLight(0x8ae8ff, 0x060711, 1.15);
      scene.add(ambient);
      const keyLight = new THREE.DirectionalLight(0xe8e8ff, 2.4);
      keyLight.position.set(-4, 8, 6);
      scene.add(keyLight);

      const waterUniforms = { uTime: { value: 0 } };
      const waterGeometry = new THREE.PlaneGeometry(42, 42, 56, 56);
      const waterMaterial = new THREE.ShaderMaterial({
        transparent: true,
        side: THREE.DoubleSide,
        uniforms: waterUniforms,
        vertexShader: `
          uniform float uTime;
          varying float vLift;
          varying vec2 vUv;
          void main() {
            vUv = uv;
            vec3 p = position;
            float broad = sin(p.x * 0.42 + uTime * 0.38) * 0.34;
            float cross = sin(p.y * 0.58 - uTime * 0.27) * 0.22;
            p.z += broad + cross;
            vLift = broad + cross;
            gl_Position = projectionMatrix * modelViewMatrix * vec4(p, 1.0);
          }
        `,
        fragmentShader: `
          varying float vLift;
          varying vec2 vUv;
          void main() {
            float rim = smoothstep(-0.34, 0.54, vLift);
            vec3 deep = vec3(0.035, 0.10, 0.18);
            vec3 light = vec3(0.38, 0.88, 0.96);
            vec3 color = mix(deep, light, rim * 0.72);
            float edge = 0.68 + sin((vUv.x + vUv.y) * 80.0) * 0.04;
            gl_FragColor = vec4(color * edge, 0.82);
          }
        `,
      });
      const water = new THREE.Mesh(waterGeometry, waterMaterial);
      water.rotation.x = -Math.PI / 2;
      water.position.y = 2.1;
      world.add(water);

      let randomState = 0x41a7c3;
      const random = () => {
        randomState = (randomState * 1664525 + 1013904223) >>> 0;
        return randomState / 4294967296;
      };

      const particleCount = window.innerWidth < 768 ? 520 : 980;
      const particlePositions = new Float32Array(particleCount * 3);
      for (let index = 0; index < particleCount; index += 1) {
        particlePositions[index * 3] = (random() - 0.5) * 22;
        particlePositions[index * 3 + 1] = 4 - random() * 39;
        particlePositions[index * 3 + 2] = 3 - random() * 18;
      }
      const particleGeometry = new THREE.BufferGeometry();
      particleGeometry.setAttribute("position", new THREE.BufferAttribute(particlePositions, 3));
      const particleMaterial = new THREE.PointsMaterial({
        color: 0x8ae8ff,
        depthWrite: false,
        opacity: 0.48,
        size: window.innerWidth < 768 ? 0.045 : 0.032,
        sizeAttenuation: true,
        transparent: true,
      });
      const particles = new THREE.Points(particleGeometry, particleMaterial);
      world.add(particles);

      const currentMaterials = [0x8ae8ff, 0x9385ff, 0x75e6ad].map((color) => new THREE.LineBasicMaterial({
        color,
        opacity: 0.24,
        transparent: true,
      }));
      for (let strand = 0; strand < 9; strand += 1) {
        const points = [];
        const phase = random() * Math.PI * 2;
        const radius = 2.2 + random() * 5.2;
        for (let step = 0; step < 18; step += 1) {
          const y = 3 - step * 2.15;
          points.push(new THREE.Vector3(
            Math.sin(step * 0.52 + phase) * radius,
            y,
            -3.5 + Math.cos(step * 0.37 + phase) * 3.2,
          ));
        }
        const curve = new THREE.CatmullRomCurve3(points, false, "centripetal", 0.5);
        const geometry = new THREE.BufferGeometry().setFromPoints(curve.getPoints(160));
        world.add(new THREE.Line(geometry, currentMaterials[strand % currentMaterials.length]));
      }

      const ark = new THREE.Group();
      const arkCoreGeometry = new THREE.IcosahedronGeometry(0.78, 1);
      const arkCoreMaterial = new THREE.MeshStandardMaterial({
        color: 0xe8e8ff,
        emissive: 0x5446c8,
        emissiveIntensity: 1.2,
        metalness: 0.34,
        roughness: 0.2,
      });
      const arkCore = new THREE.Mesh(arkCoreGeometry, arkCoreMaterial);
      ark.add(arkCore);

      const arkShell = new THREE.Mesh(
        new THREE.IcosahedronGeometry(1.08, 1),
        new THREE.MeshBasicMaterial({ color: 0x8ae8ff, opacity: 0.34, transparent: true, wireframe: true }),
      );
      ark.add(arkShell);
      const arkRing = new THREE.Mesh(
        new THREE.TorusGeometry(1.32, 0.018, 6, 96),
        new THREE.MeshBasicMaterial({ color: 0x9385ff, opacity: 0.72, transparent: true }),
      );
      arkRing.rotation.x = Math.PI / 2.7;
      ark.add(arkRing);

      const shardGeometry = new THREE.OctahedronGeometry(0.12, 0);
      for (let index = 0; index < 6; index += 1) {
        const available = index < 4;
        const shard = new THREE.Mesh(
          shardGeometry,
          new THREE.MeshBasicMaterial({
            color: available ? 0x8ae8ff : 0x596071,
            opacity: available ? 0.96 : 0.28,
            transparent: true,
            wireframe: !available,
          }),
        );
        const angle = (index / 6) * Math.PI * 2;
        shard.position.set(Math.cos(angle) * 1.6, Math.sin(angle) * 1.6, 0);
        ark.add(shard);
      }
      const arkLight = new THREE.PointLight(0x8ae8ff, 5.5, 11, 1.8);
      ark.add(arkLight);
      world.add(ark);

      const makeMarker = (color: number, position: [number, number, number], scale: number) => {
        const marker = new THREE.Group();
        marker.position.set(...position);
        const orb = new THREE.Mesh(
          new THREE.SphereGeometry(scale, 16, 16),
          new THREE.MeshBasicMaterial({ color }),
        );
        marker.add(orb);
        const ring = new THREE.Mesh(
          new THREE.TorusGeometry(scale * 2.2, scale * 0.08, 6, 48),
          new THREE.MeshBasicMaterial({ color, opacity: 0.52, transparent: true }),
        );
        ring.rotation.x = Math.PI / 2;
        marker.add(ring);
        world.add(marker);
        return marker;
      };

      const lostOrigin = makeMarker(0xffc977, [-2.6, -4.8, -1.8], 0.18);
      const freshNode = makeMarker(0x75e6ad, [2.8, -19.2, -2.6], 0.24);
      const authority = makeMarker(0xffc977, [0, -30.5, -2.5], 0.34);
      authority.add(new THREE.PointLight(0xffc977, 7, 14, 1.6));

      const cameraKeys = [
        { position: [0, 3.6, 10.5], target: [0, 1.1, 0] },
        { position: [1.8, -2.7, 8.3], target: [-0.4, -5.3, -1.4] },
        { position: [-1.6, -9.8, 7.2], target: [0.2, -12.8, -2.0] },
        { position: [1.4, -17.7, 7.8], target: [0, -20.7, -2.1] },
        { position: [0, -27.1, 7.0], target: [0, -30.0, -2.4] },
      ];
      const positionCurve = new THREE.CatmullRomCurve3(
        cameraKeys.map((key) => new THREE.Vector3(...key.position)), false, "centripetal", 0.5,
      );
      const targetCurve = new THREE.CatmullRomCurve3(
        cameraKeys.map((key) => new THREE.Vector3(...key.target)), false, "centripetal", 0.5,
      );
      const samples = 420;
      const cameraPositions = Array.from({ length: samples }, (_, index) => positionCurve.getPointAt(index / (samples - 1)));
      const cameraTargets = Array.from({ length: samples }, (_, index) => targetCurve.getPointAt(index / (samples - 1)));

      const timer = new THREE.Timer();
      timer.connect(document);
      const pointer = new THREE.Vector2();
      const pointerTarget = new THREE.Vector2();
      const lookTarget = new THREE.Vector3();
      let smoothProgress = progressRef.current;
      let visible = true;

      const resize = () => {
        const { width, height } = host.getBoundingClientRect();
        if (!width || !height) {
          return;
        }
        renderer.setSize(width, height, false);
        renderer.setPixelRatio(Math.min(window.devicePixelRatio, width < 768 ? 1 : 1.5));
        camera.aspect = width / height;
        camera.updateProjectionMatrix();
      };

      const onPointerMove = (event: PointerEvent) => {
        pointerTarget.set(
          (event.clientX / window.innerWidth - 0.5) * 2,
          (event.clientY / window.innerHeight - 0.5) * 2,
        );
      };

      const render = (timestamp: number) => {
        timer.update(timestamp);
        const delta = Math.min(timer.getDelta(), 0.05);
        const elapsed = timer.getElapsed();
        const smoothing = 1 - Math.exp(-5.5 * delta);
        smoothProgress += (progressRef.current - smoothProgress) * smoothing;
        pointer.lerp(pointerTarget, 1 - Math.exp(-3.2 * delta));

        const sampleFloat = smoothProgress * (samples - 1);
        const sampleA = Math.floor(sampleFloat);
        const sampleB = Math.min(sampleA + 1, samples - 1);
        const fraction = sampleFloat - sampleA;
        camera.position.lerpVectors(cameraPositions[sampleA], cameraPositions[sampleB], fraction);
        camera.position.x += pointer.x * 0.28;
        camera.position.y -= pointer.y * 0.12;
        lookTarget.lerpVectors(cameraTargets[sampleA], cameraTargets[sampleB], fraction);
        lookTarget.x += pointer.x * 0.2;
        camera.lookAt(lookTarget);

        ark.position.set(
          Math.sin(smoothProgress * Math.PI * 2.2) * 1.15,
          -5.8 - smoothProgress * 23.0,
          -2.2 + Math.cos(smoothProgress * Math.PI * 1.4) * 0.45,
        );
        ark.rotation.y = elapsed * 0.17 + smoothProgress * Math.PI * 1.8;
        ark.rotation.x = Math.sin(elapsed * 0.28) * 0.12;
        arkRing.rotation.z = elapsed * 0.24;
        ark.children.slice(3, 9).forEach((child, index) => {
          child.rotation.x = elapsed * (0.5 + index * 0.04);
          child.rotation.y = elapsed * 0.4;
        });

        waterUniforms.uTime.value = elapsed;
        waterMaterial.opacity = Math.max(0, 1 - smoothProgress * 3.4);
        water.visible = waterMaterial.opacity > 0.02;
        particles.rotation.y = Math.sin(elapsed * 0.035) * 0.08;
        lostOrigin.rotation.z = elapsed * 0.2;
        freshNode.rotation.z = -elapsed * 0.17;
        authority.rotation.z = elapsed * 0.11;

        currentColor.lerpColors(surfaceColor, abyssColor, Math.min(1, smoothProgress * 1.12));
        renderer.setClearColor(currentColor, 1);
        fog.color.copy(currentColor);
        fog.density = 0.038 + smoothProgress * 0.017;
        renderer.render(scene, camera);
      };

      const syncLoop = () => renderer.setAnimationLoop(visible && !motionQuery.matches ? render : null);
      const onMotionChange = () => syncLoop();
      const resizeObserver = new ResizeObserver(resize);
      const intersectionObserver = new IntersectionObserver(([entry]) => {
        visible = entry.isIntersecting;
        syncLoop();
      }, { rootMargin: "100px" });

      resizeObserver.observe(host);
      intersectionObserver.observe(section);
      window.addEventListener("pointermove", onPointerMove, { passive: true });
      motionQuery.addEventListener("change", onMotionChange);
      resize();
      syncLoop();

      disposeScene = () => {
        renderer.setAnimationLoop(null);
        resizeObserver.disconnect();
        intersectionObserver.disconnect();
        window.removeEventListener("pointermove", onPointerMove);
        motionQuery.removeEventListener("change", onMotionChange);
        timer.dispose();
        scene.traverse((object) => {
          if (object instanceof THREE.Mesh || object instanceof THREE.Line || object instanceof THREE.Points) {
            object.geometry.dispose();
            const materials = Array.isArray(object.material) ? object.material : [object.material];
            materials.forEach((material) => material.dispose());
          }
        });
        renderer.dispose();
        renderer.domElement.remove();
      };
    }

    let sceneStarted = false;
    const beginJourney = () => {
      if (sceneStarted) {
        return;
      }
      sceneStarted = true;
      window.removeEventListener("scroll", beginJourney);
      mountScene().catch(() => {
        host.dataset.renderState = "fallback";
      });
    };

    if (window.scrollY > 4) {
      beginJourney();
    } else {
      window.addEventListener("scroll", beginJourney, { passive: true, once: true });
    }

    return () => {
      cancelled = true;
      window.removeEventListener("scroll", beginJourney);
      disposeScene?.();
    };
  }, []);

  return (
    <section ref={sectionRef} className="ocean-journey" aria-label="ArkMesh continuity descent">
      <div className="ocean-stage">
        <div ref={canvasHostRef} className="ocean-canvas" aria-hidden="true" />
        <div className="ocean-vignette" aria-hidden="true" />
        <div className="ocean-caustics" aria-hidden="true" />

        <div className="ocean-hud" aria-hidden="true">
          <div className="ocean-hud-brand"><span>AM</span><strong>ARKMESH</strong><small>CONTINUITY DIVE / 001</small></div>
          <div className="ocean-hud-status"><i /> LINK LOCAL <span>PEERS PLANNED</span></div>
        </div>

        <a className="ocean-skip" href="#continuity-visual">Skip journey</a>

        <div className="ocean-depth" aria-hidden="true">
          <span ref={depthRef}>0000 M</span>
          <div><i /></div>
          <small>DEPTH<br />BELOW ORIGIN</small>
        </div>

        <div className="ocean-reticle" aria-hidden="true"><i /><span>ARK CAPSULE</span></div>

        <div className="ocean-chapters">
          {chapters.map((chapter, index) => (
            <article
              key={chapter.eyebrow}
              className={`ocean-chapter${activeChapter === index ? " is-active" : ""}`}
              aria-hidden={activeChapter !== index}
            >
              <div className="ocean-chapter-index"><span>{chapter.eyebrow}</span><strong>{chapter.depth}</strong></div>
              {index === 0 ? <h1>{chapter.title}</h1> : <h2>{chapter.title}</h2>}
              <p>{chapter.body}</p>
              {index === 0 && (
                <div className="ocean-dive-cue" aria-hidden="true"><i /><span>Scroll to descend</span></div>
              )}
              {index === chapters.length - 1 && (
                <div className="ocean-actions">
                  <a href="/wiki#overview">Enter the archive</a>
                  <a href="https://github.com/AlleyBo55/arkmesh/issues/new?template=research-proposal.yml" target="_blank" rel="noreferrer">Join the research ↗</a>
                </div>
              )}
            </article>
          ))}
        </div>

        <div className="ocean-principles" aria-hidden="true">
          <span>AUTHENTICITY BEFORE AVAILABILITY</span>
          <span>NO CENTRAL TRACKER</span>
          <span>HUMAN CONSENT REQUIRED</span>
        </div>

      </div>

      <div className="ocean-reduced-story">
        {chapters.map((chapter, index) => (
          <article key={chapter.eyebrow}>
            <span>{chapter.eyebrow} · {chapter.depth}</span>
            {index === 0 ? <h1>{chapter.title}</h1> : <h2>{chapter.title}</h2>}
            <p>{chapter.body}</p>
          </article>
        ))}
      </div>
    </section>
  );
}
