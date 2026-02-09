import { JSX } from "preact";

type LabelProps = {
    children: JSX.Element | string;
    htmlFor?: string;
    class?: string;
};

export function Label({ children, htmlFor, class: className = "" }: LabelProps) {
    return (
        <label
            htmlFor={htmlFor}
            class={`text-sm font-medium leading-none text-zinc-300 peer-disabled:cursor-not-allowed peer-disabled:opacity-70 ${className}`}>
            {children}
        </label>
    );
}
