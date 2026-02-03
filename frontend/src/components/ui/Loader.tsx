import { useEffect, useState } from "preact/hooks";
import { HTMLAttributes } from "preact";

type Variant = "noise" | "scanline";
type Size = "xs" | "sm" | "md" | "lg";

type LoaderProps = HTMLAttributes<HTMLDivElement> & {
    variant?: Variant;
    size?: Size;
};

const CHARS = ["█", "▓", "▒", "░", "■", "□", "•", " "];

const COLS = 6;
const ROWS = 3;

function randomChar(chars: string[]) {
    return chars[Math.floor(Math.random() * chars.length)];
}

function generateGrid(chars: string[]) {
    return Array.from({ length: ROWS }, () =>
        Array.from({ length: COLS }, () => randomChar(chars)),
    );
}

const base = "inline-block font-bold whitespace-pre leading-none";

const sizes: Record<Size, string> = {
    xs: "text-xs",
    sm: "text-sm",
    md: "text-base",
    lg: "text-lg",
};

export function Loader({
    variant = "noise",
    size = "md",
    class: className,
    ...props
}: LoaderProps) {
    const [grid, setGrid] = useState(() => generateGrid(CHARS));
    const [blink, setBlink] = useState(true);

    useEffect(() => {
        const id = setInterval(() => {
            setGrid((prev) =>
                prev.map((row) =>
                    row.map(() =>
                        Math.random() > 0.6 ?
                            randomChar(CHARS)
                        :   row[Math.floor(Math.random() * row.length)],
                    ),
                ),
            );
        }, 120);

        return () => clearInterval(id);
    }, []);

    useEffect(() => {
        const id = setInterval(() => {
            if (variant === "scanline") {
                setBlink(!blink);
            }
        }, 500);
        return () => clearInterval(id);
    }, [blink]);

    if (variant === "scanline") {
        return (
            <div {...props} class={[base, sizes[size], className].filter(Boolean).join(" ")}>
                {blink ? " █" : " _"}
            </div>
        );
    }

    return (
        <div
            {...props}
            class={["text-lime-400", base, sizes[size], className].filter(Boolean).join(" ")}>
            {grid.map((row) => (
                <div>{row.join("")}</div>
            ))}
        </div>
    );
}
