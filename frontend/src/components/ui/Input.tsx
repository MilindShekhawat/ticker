import { InputHTMLAttributes } from "preact";

type InputProps = InputHTMLAttributes & {
    error?: string;
};

const base = `inline-flex w-full h-8 px-4 text-sm mt-1 font-medium border bg-zinc-900 transition outline-none disabled:pointer-events-none disabled:opacity-40`;
const primary = `placeholder:text-zinc-500  text-zinc-100   border-zinc-700   focus:border-lime-400   focus:text-lime-400`;
const error_ = `placeholder:text-orange-800 text-orange-600 border-orange-600 focus:border-orange-600 focus:text-orange-600`;

export function Input({ class: className, error, ...props }: InputProps) {
    return (
        <div class="w-full">
            <input
                {...props}
                class={[base, error ? error_ : primary, className].filter(Boolean).join(" ")}
            />
            {error && <p class="mt-1 text-sm text-orange-600">{error}</p>}
        </div>
    );
}
