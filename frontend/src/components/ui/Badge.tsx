import { ComponentChildren, HTMLAttributes } from "preact";

type Variant = "primary" | "success" | "destructive" | "warning" | "blue" | "purple";

type BadgeProps = HTMLAttributes<HTMLDivElement> & {
    children: ComponentChildren;
    variant?: Variant;
};

const base =
    "inline-flex items-center border-transparent px-1.5 py-0.5 text-xs font-bold transition-colors";

const variants: Record<Variant, string> = {
    primary: "bg-lime-400/40 text-lime-400",
    success: "bg-green-500/40 text-green-400",
    destructive: "bg-orange-600/40 text-orange-500",
    warning: "bg-yellow-500/40 text-yellow-400",
    blue: "bg-sky-500/40 text-sky-400",
    purple: "bg-purple-500/40 text-purple-400",
};

export function Badge({ children, variant = "primary", class: className, ...props }: BadgeProps) {
    return (
        <div {...props} class={[base, variants[variant], className].filter(Boolean).join(" ")}>
            {children}
        </div>
    );
}
