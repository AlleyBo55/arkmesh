export function HeroMesh() {
  return (
    <figure
      className="hero-mesh"
      role="img"
      aria-label="ArkMesh continuity map showing a lost publisher, four surviving shards, a fresh node, and required human authority"
    >
      <div className="hero-mesh-aura" aria-hidden="true" />
      <svg className="hero-mesh-svg" viewBox="0 0 640 640" aria-hidden="true">
        <defs>
          <linearGradient id="meshBeam" x1="0" x2="1">
            <stop offset="0" stopColor="#65f1ff" stopOpacity="0" />
            <stop offset="0.48" stopColor="#65f1ff" stopOpacity="1" />
            <stop offset="1" stopColor="#65f1ff" stopOpacity="0" />
          </linearGradient>
          <radialGradient id="meshCore" cx="50%" cy="50%" r="50%">
            <stop offset="0" stopColor="#65f1ff" stopOpacity="0.34" />
            <stop offset="1" stopColor="#65f1ff" stopOpacity="0" />
          </radialGradient>
        </defs>

        <circle cx="320" cy="320" r="266" className="mesh-orbit mesh-orbit-outer" />
        <circle cx="320" cy="320" r="205" className="mesh-orbit mesh-orbit-middle" />
        <circle cx="320" cy="320" r="142" className="mesh-orbit mesh-orbit-inner" />
        <circle cx="320" cy="320" r="118" fill="url(#meshCore)" />

        <g className="mesh-beams">
          <path d="M320 320 L130 130" />
          <path d="M320 320 L510 130" />
          <path d="M320 320 L535 350" />
          <path d="M320 320 L450 535" />
          <path d="M320 320 L190 535" />
          <path d="M320 320 L105 350" />
        </g>

        <g className="mesh-packets">
          <path d="M320 320 L130 130" pathLength="1" />
          <path d="M320 320 L510 130" pathLength="1" />
          <path d="M320 320 L535 350" pathLength="1" />
          <path d="M320 320 L450 535" pathLength="1" />
          <path d="M320 320 L190 535" pathLength="1" />
          <path d="M320 320 L105 350" pathLength="1" />
        </g>

        <g className="mesh-nodes">
          <circle cx="130" cy="130" r="13" className="mesh-node-lost" />
          <circle cx="510" cy="130" r="13" />
          <circle cx="535" cy="350" r="13" />
          <circle cx="450" cy="535" r="13" />
          <circle cx="190" cy="535" r="13" />
          <circle cx="105" cy="350" r="13" />
        </g>
      </svg>

      <div className="hero-core" aria-hidden="true">
        <span className="hero-core-index">AM://01</span>
        <strong>ARK<br />MESH</strong>
        <span className="hero-core-state">RECOVERABLE</span>
      </div>

      <div className="mesh-label mesh-label-one">
        <span>ORIGIN</span>
        <strong>LOST</strong>
      </div>
      <div className="mesh-label mesh-label-two">
        <span>SHARDS</span>
        <strong>04 / 06</strong>
      </div>
      <div className="mesh-label mesh-label-three">
        <span>FRESH NODE</span>
        <strong>WAITING</strong>
      </div>
      <div className="mesh-label mesh-label-four">
        <span>AUTHORITY</span>
        <strong>HUMAN</strong>
      </div>

      <figcaption className="hero-mesh-caption">
        <span className="hero-mesh-live">LIVE THESIS</span>
        <span>NO SINGLE POINT OF FAILURE</span>
        <span>PEER NETWORK: PLANNED</span>
      </figcaption>
    </figure>
  );
}
