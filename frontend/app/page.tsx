import Campus from "@/components/map/campus";

export default function Home() {
    return (
        <main>
            <Campus />
            <div
                style={{
                    position: "absolute",
                    top: 20,
                    left: 20,
                    zIndex: 10,
                    color: "black",
                    background: "rgba(255,255,255,0.8)",
                    padding: "1rem",
                    borderRadius: "8px",
                }}
            >
                <h1 style={{ margin: 0 }}>Ghost Map</h1>
                <p style={{ margin: 0 }}>GMU Campus</p>
            </div>
        </main>
    );
}
