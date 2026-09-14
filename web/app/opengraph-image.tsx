import { ImageResponse } from "next/og";

export const alt = "ArkMesh: verifiable offline AI capability preservation";
export const size = { width: 1200, height: 630 };
export const contentType = "image/png";

export default function OpenGraphImage() {
  return new ImageResponse(
    (
      <div
        style={{
          alignItems: "stretch",
          background: "#00072f",
          color: "#ecfaff",
          display: "flex",
          flexDirection: "column",
          fontFamily: "monospace",
          height: "100%",
          justifyContent: "space-between",
          padding: "54px",
          width: "100%",
        }}
      >
        <div style={{ alignItems: "center", display: "flex", justifyContent: "space-between" }}>
          <div style={{ border: "4px solid #65f1ff", color: "#65f1ff", display: "flex", fontSize: 30, fontWeight: 800, padding: "12px 18px" }}>ARKMESH</div>
          <div style={{ color: "#8dbce7", display: "flex", fontSize: 24 }}>v0.1.0-alpha.1 // RESEARCH PROTOTYPE</div>
        </div>
        <div style={{ display: "flex", flexDirection: "column" }}>
          <div style={{ color: "#65f1ff", display: "flex", fontSize: 25, letterSpacing: "0.16em" }}>OPEN-SOURCE AI CONTINUITY RESEARCH</div>
          <div style={{ display: "flex", fontSize: 86, fontWeight: 900, lineHeight: 1.02, marginTop: 24, maxWidth: 1040 }}>What if AI<br />goes dark?</div>
        </div>
        <div style={{ borderTop: "3px solid #69c9ff", color: "#8dbce7", display: "flex", fontSize: 24, justifyContent: "space-between", paddingTop: 24 }}>
          <span>CONTENT ADDRESSED</span><span>DONORLESS RECOVERY</span><span>EXPLICIT TRUST</span>
        </div>
      </div>
    ),
    size,
  );
}
