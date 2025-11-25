import { cn } from "@/lib/utils";

describe("Utility functions", () => {
  describe("cn (className merger)", () => {
    it("should merge class names correctly", () => {
      const result = cn("base-class", "additional-class");
      expect(result).toContain("base-class");
      expect(result).toContain("additional-class");
    });

    it("should handle conditional classes", () => {
      const result = cn("base", true && "conditional", false &&" skipped");
      expect(result).toContain("base");
      expect(result).toContain("conditional");
      expect(result).not.toContain("skipped");
    });

    it("should merge Tailwind classes correctly", () => {
      // When the same Tailwind property is specified multiple times,
      // tailwind-merge should keep only the last one
      const result = cn("px-2", "px-4");
      expect(result).toBe("px-4");
    });

    it("should handle undefined and null values", () => {
      const result = cn("base", undefined, null, "last");
      expect(result).toContain("base");
      expect(result).toContain("last");
      expect(result).not.toContain("undefined");
      expect(result).not.toContain("null");
    });

    it("should handle arrays of classes", () => {
      const result = cn(["class1", "class2"], "class3");
      expect(result).toContain("class1");
      expect(result).toContain("class2");
      expect(result).toContain("class3");
    });

    it("should handle empty inputs", () => {
      const result = cn();
      expect(result).toBe("");
    });
  });
});
