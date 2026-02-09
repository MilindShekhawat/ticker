import { useEffect, useState } from "preact/hooks";
import { HTMLAttributes } from "preact";
import { cva, type VariantProps } from "class-variance-authority";

const loaderVariants = cva(`inline-block font-bold whitespace-pre leading-none`, {
    variants: {
        variant: {
            noise: `text-lime-400`,
            scanline: ``,
        },
        size: {
            xs: `text-xs`,
            sm: `text-sm`,
            md: `text-base`,
            lg: `text-lg`,
        },
    },
    defaultVariants: {
        variant: `noise`,
        size: `md`,
    },
});

type LoaderProps = HTMLAttributes<HTMLDivElement> & VariantProps<typeof loaderVariants>;

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

export function Loader({ variant, size, class: className, ...props }: LoaderProps) {
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
                setBlink((b) => !b);
            }
        }, 400);

        return () => clearInterval(id);
    }, [variant]);

    if (variant === "scanline") {
        return (
            <div {...props} class={loaderVariants({ variant, size, class: className })}>
                {blink ? " █" : " _"}
            </div>
        );
    }

    return (
        <div {...props} class={loaderVariants({ variant, size, class: className })}>
            {grid.map((row) => (
                <div>{row.join("")}</div>
            ))}
        </div>
    );
}
