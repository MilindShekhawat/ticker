import { useEffect, useMemo, useRef, useState } from "preact/hooks";
import { createPortal } from "preact/compat";
import { cva, type VariantProps } from "class-variance-authority";

const triggerVariants = cva(
    `inline-flex w-full h-8 px-3 pr-8 text-sm font-medium border bg-zinc-900 transition outline-none disabled:pointer-events-none disabled:opacity-40 items-center justify-between`,
    {
        variants: {
            state: {
                default: `text-zinc-100 border-zinc-700 focus:border-lime-400 focus:text-lime-400`,
                error: `text-orange-600 border-orange-600 focus:border-orange-600 focus:text-orange-600`,
            },
        },
        defaultVariants: {
            state: `default`,
        },
    },
);

type SelectOption = {
    value: string;
    label: string;
};

type SelectProps = {
    value: string;
    options: SelectOption[];
    onChange: (value: string) => void;
    class?: string;
    error?: string;
    disabled?: boolean;
} & VariantProps<typeof triggerVariants>;

export function Select({
    value,
    options,
    onChange,
    class: className,
    error,
    disabled,
    state,
}: SelectProps) {
    const rootRef = useRef<HTMLDivElement>(null);
    const triggerRef = useRef<HTMLButtonElement>(null);
    const listRef = useRef<HTMLUListElement>(null);
    const [open, setOpen] = useState(false);
    const [activeIndex, setActiveIndex] = useState(-1);
    const [panelStyle, setPanelStyle] = useState<Record<string, string>>({});

    const selectedIndex = useMemo(
        () => options.findIndex((option) => option.value === value),
        [options, value],
    );
    const selectedLabel =
        selectedIndex >= 0 ? options[selectedIndex].label : (options[0]?.label ?? "Select");

    useEffect(() => {
        if (!open) return;

        const onPointerDown = (event: MouseEvent) => {
            if (!rootRef.current?.contains(event.target as Node)) {
                setOpen(false);
            }
        };

        const onEscape = (event: KeyboardEvent) => {
            if (event.key === "Escape") {
                setOpen(false);
            }
        };

        document.addEventListener("mousedown", onPointerDown);
        document.addEventListener("keydown", onEscape);
        return () => {
            document.removeEventListener("mousedown", onPointerDown);
            document.removeEventListener("keydown", onEscape);
        };
    }, [open]);

    useEffect(() => {
        if (!open) return;
        const indexToFocus = selectedIndex >= 0 ? selectedIndex : 0;
        setActiveIndex(indexToFocus);
    }, [open, selectedIndex]);

    useEffect(() => {
        if (!open || activeIndex < 0) return;
        const optionEl = listRef.current?.querySelector<HTMLLIElement>(
            `[data-index="${activeIndex}"]`,
        );
        optionEl?.scrollIntoView({ block: "nearest" });
    }, [open, activeIndex]);

    useEffect(() => {
        if (!open) return;

        const updatePanelPosition = () => {
            const rect = triggerRef.current?.getBoundingClientRect();
            if (!rect) return;
            setPanelStyle({
                position: "fixed",
                top: `${rect.bottom + 4}px`,
                left: `${rect.left}px`,
                width: `${rect.width}px`,
                zIndex: "80",
            });
        };

        updatePanelPosition();
        window.addEventListener("resize", updatePanelPosition);
        window.addEventListener("scroll", updatePanelPosition, true);
        return () => {
            window.removeEventListener("resize", updatePanelPosition);
            window.removeEventListener("scroll", updatePanelPosition, true);
        };
    }, [open]);

    const handleTriggerKeyDown = (event: KeyboardEvent) => {
        if (disabled) return;

        if (event.key === "Enter" || event.key === " ") {
            event.preventDefault();
            setOpen((prev) => !prev);
            return;
        }
        if (event.key === "ArrowDown") {
            event.preventDefault();
            setOpen(true);
            setActiveIndex((prev) => {
                if (prev < 0) return selectedIndex >= 0 ? selectedIndex : 0;
                return Math.min(prev + 1, options.length - 1);
            });
            return;
        }
        if (event.key === "ArrowUp") {
            event.preventDefault();
            setOpen(true);
            setActiveIndex((prev) => {
                if (prev < 0) return selectedIndex >= 0 ? selectedIndex : options.length - 1;
                return Math.max(prev - 1, 0);
            });
            return;
        }
        if (event.key === "Escape") {
            setOpen(false);
        }
    };

    const handleListKeyDown = (event: KeyboardEvent) => {
        if (event.key === "ArrowDown") {
            event.preventDefault();
            setActiveIndex((prev) => Math.min(prev + 1, options.length - 1));
            return;
        }
        if (event.key === "ArrowUp") {
            event.preventDefault();
            setActiveIndex((prev) => Math.max(prev - 1, 0));
            return;
        }
        if (event.key === "Enter" || event.key === " ") {
            event.preventDefault();
            if (activeIndex >= 0) {
                onChange(options[activeIndex].value);
                setOpen(false);
            }
            return;
        }
        if (event.key === "Escape") {
            event.preventDefault();
            setOpen(false);
        }
    };

    return (
        <div class="w-full relative" ref={rootRef}>
            <button
                ref={triggerRef}
                type="button"
                disabled={disabled}
                aria-haspopup="listbox"
                aria-expanded={open}
                class={triggerVariants({
                    state: error ? "error" : state,
                    class: className,
                })}
                onClick={() => !disabled && setOpen((prev) => !prev)}
                onKeyDown={handleTriggerKeyDown}>
                <span class="truncate">{selectedLabel}</span>
                <svg
                    aria-hidden="true"
                    viewBox="0 0 16 16"
                    class={`absolute right-2 top-1/2 h-3.5 w-3.5 -translate-y-1/2 text-zinc-500 transition ${
                        open ? "rotate-180" : ""
                    }`}>
                    <path
                        d="M4.25 6.25 8 10l3.75-3.75"
                        fill="none"
                        stroke="currentColor"
                        stroke-linecap="round"
                        stroke-linejoin="round"
                        stroke-width="1.5"
                    />
                </svg>
            </button>

            {open &&
                createPortal(
                    <ul
                        ref={listRef}
                        role="listbox"
                        tabIndex={-1}
                        onKeyDown={handleListKeyDown}
                        style={panelStyle}
                        class="max-h-60 overflow-auto border border-zinc-700 bg-zinc-900 shadow-lg">
                        {options.map((option, index) => {
                            const isSelected = option.value === value;
                            const isActive = index === activeIndex;
                            return (
                                <li
                                    key={option.value}
                                    data-index={index}
                                    role="option"
                                    aria-selected={isSelected}
                                    class={`cursor-pointer px-3 py-2 text-xs uppercase tracking-wider ${
                                        isSelected ? "bg-lime-400 text-zinc-950"
                                        : isActive ? "bg-zinc-800 text-zinc-100"
                                        : "bg-zinc-900 text-zinc-100 hover:bg-zinc-800"
                                    }`}
                                    onMouseEnter={() => setActiveIndex(index)}
                                    onMouseDown={(event) => {
                                        event.preventDefault();
                                        onChange(option.value);
                                        setOpen(false);
                                    }}>
                                    {option.label}
                                </li>
                            );
                        })}
                    </ul>,
                    document.body,
                )}

            {error && <p class="mt-1 text-sm text-orange-600">{error}</p>}
        </div>
    );
}
