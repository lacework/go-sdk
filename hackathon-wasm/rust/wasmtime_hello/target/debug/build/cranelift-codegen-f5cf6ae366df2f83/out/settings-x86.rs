#[derive(Clone, Hash)]
/// Flags group `x86`.
pub struct Flags {
    bytes: [u8; 5],
}
impl Flags {
    /// Create flags x86 settings group.
    #[allow(unused_variables)]
    pub fn new(shared: &settings::Flags, builder: Builder) -> Self {
        let bvec = builder.state_for("x86");
        let mut x86 = Self { bytes: [0; 5] };
        debug_assert_eq!(bvec.len(), 2);
        x86.bytes[0..2].copy_from_slice(&bvec);
        // Precompute #16.
        if shared.enable_simd() && x86.has_avx2() {
            x86.bytes[2] |= 1 << 0;
        }
        // Precompute #17.
        if shared.enable_simd() && x86.has_avx512bitalg() {
            x86.bytes[2] |= 1 << 1;
        }
        // Precompute #18.
        if shared.enable_simd() && x86.has_avx512dq() {
            x86.bytes[2] |= 1 << 2;
        }
        // Precompute #19.
        if shared.enable_simd() && x86.has_avx512f() {
            x86.bytes[2] |= 1 << 3;
        }
        // Precompute #20.
        if shared.enable_simd() && x86.has_avx512vbmi() {
            x86.bytes[2] |= 1 << 4;
        }
        // Precompute #21.
        if shared.enable_simd() && x86.has_avx512vl() {
            x86.bytes[2] |= 1 << 5;
        }
        // Precompute #22.
        if shared.enable_simd() && x86.has_avx() {
            x86.bytes[2] |= 1 << 6;
        }
        // Precompute #23.
        if x86.has_bmi1() {
            x86.bytes[2] |= 1 << 7;
        }
        // Precompute #24.
        if x86.has_avx() && x86.has_fma() {
            x86.bytes[3] |= 1 << 0;
        }
        // Precompute #25.
        if x86.has_lzcnt() {
            x86.bytes[3] |= 1 << 1;
        }
        // Precompute #26.
        if x86.has_popcnt() && x86.has_sse42() {
            x86.bytes[3] |= 1 << 2;
        }
        // Precompute #27.
        if x86.has_sse41() {
            x86.bytes[3] |= 1 << 3;
        }
        // Precompute #28.
        if shared.enable_simd() && x86.has_sse41() {
            x86.bytes[3] |= 1 << 4;
        }
        // Precompute #29.
        if x86.has_sse41() && x86.has_sse42() {
            x86.bytes[3] |= 1 << 5;
        }
        // Precompute #30.
        if shared.enable_simd() && x86.has_sse41() && x86.has_sse42() {
            x86.bytes[3] |= 1 << 6;
        }
        // Precompute #31.
        if x86.has_ssse3() {
            x86.bytes[3] |= 1 << 7;
        }
        // Precompute #32.
        if shared.enable_simd() && x86.has_ssse3() {
            x86.bytes[4] |= 1 << 0;
        }
        x86
    }
}
impl Flags {
    /// Iterates the setting values.
    pub fn iter(&self) -> impl Iterator<Item = Value> {
        let mut bytes = [0; 2];
        bytes.copy_from_slice(&self.bytes[0..2]);
        DESCRIPTORS.iter().filter_map(move |d| {
            let values = match &d.detail {
                detail::Detail::Preset => return None,
                detail::Detail::Enum { last, enumerators } => Some(TEMPLATE.enums(*last, *enumerators)),
                _ => None
            };
            Some(Value{ name: d.name, detail: d.detail, values, value: bytes[d.offset as usize] })
        })
    }
}
/// User-defined settings.
#[allow(dead_code)]
impl Flags {
    /// Get a view of the boolean predicates.
    pub fn predicate_view(&self) -> crate::settings::PredicateView {
        crate::settings::PredicateView::new(&self.bytes[0..])
    }
    /// Dynamic numbered predicate getter.
    fn numbered_predicate(&self, p: usize) -> bool {
        self.bytes[0 + p / 8] & (1 << (p % 8)) != 0
    }
    /// Has support for SSE3.
    /// SSE3: CPUID.01H:ECX.SSE3[bit 0]
    pub fn has_sse3(&self) -> bool {
        self.numbered_predicate(0)
    }
    /// Has support for SSSE3.
    /// SSSE3: CPUID.01H:ECX.SSSE3[bit 9]
    pub fn has_ssse3(&self) -> bool {
        self.numbered_predicate(1)
    }
    /// Has support for SSE4.1.
    /// SSE4.1: CPUID.01H:ECX.SSE4_1[bit 19]
    pub fn has_sse41(&self) -> bool {
        self.numbered_predicate(2)
    }
    /// Has support for SSE4.2.
    /// SSE4.2: CPUID.01H:ECX.SSE4_2[bit 20]
    pub fn has_sse42(&self) -> bool {
        self.numbered_predicate(3)
    }
    /// Has support for AVX.
    /// AVX: CPUID.01H:ECX.AVX[bit 28]
    pub fn has_avx(&self) -> bool {
        self.numbered_predicate(4)
    }
    /// Has support for AVX2.
    /// AVX2: CPUID.07H:EBX.AVX2[bit 5]
    pub fn has_avx2(&self) -> bool {
        self.numbered_predicate(5)
    }
    /// Has support for FMA.
    /// FMA: CPUID.01H:ECX.FMA[bit 12]
    pub fn has_fma(&self) -> bool {
        self.numbered_predicate(6)
    }
    /// Has support for AVX512BITALG.
    /// AVX512BITALG: CPUID.07H:ECX.AVX512BITALG[bit 12]
    pub fn has_avx512bitalg(&self) -> bool {
        self.numbered_predicate(7)
    }
    /// Has support for AVX512DQ.
    /// AVX512DQ: CPUID.07H:EBX.AVX512DQ[bit 17]
    pub fn has_avx512dq(&self) -> bool {
        self.numbered_predicate(8)
    }
    /// Has support for AVX512VL.
    /// AVX512VL: CPUID.07H:EBX.AVX512VL[bit 31]
    pub fn has_avx512vl(&self) -> bool {
        self.numbered_predicate(9)
    }
    /// Has support for AVX512VMBI.
    /// AVX512VBMI: CPUID.07H:ECX.AVX512VBMI[bit 1]
    pub fn has_avx512vbmi(&self) -> bool {
        self.numbered_predicate(10)
    }
    /// Has support for AVX512F.
    /// AVX512F: CPUID.07H:EBX.AVX512F[bit 16]
    pub fn has_avx512f(&self) -> bool {
        self.numbered_predicate(11)
    }
    /// Has support for POPCNT.
    /// POPCNT: CPUID.01H:ECX.POPCNT[bit 23]
    pub fn has_popcnt(&self) -> bool {
        self.numbered_predicate(12)
    }
    /// Has support for BMI1.
    /// BMI1: CPUID.(EAX=07H, ECX=0H):EBX.BMI1[bit 3]
    pub fn has_bmi1(&self) -> bool {
        self.numbered_predicate(13)
    }
    /// Has support for BMI2.
    /// BMI2: CPUID.(EAX=07H, ECX=0H):EBX.BMI2[bit 8]
    pub fn has_bmi2(&self) -> bool {
        self.numbered_predicate(14)
    }
    /// Has support for LZCNT.
    /// LZCNT: CPUID.EAX=80000001H:ECX.LZCNT[bit 5]
    pub fn has_lzcnt(&self) -> bool {
        self.numbered_predicate(15)
    }
    /// Computed predicate `shared.enable_simd() && x86.has_avx2()`.
    pub fn use_avx2_simd(&self) -> bool {
        self.numbered_predicate(16)
    }
    /// Computed predicate `shared.enable_simd() && x86.has_avx512bitalg()`.
    pub fn use_avx512bitalg_simd(&self) -> bool {
        self.numbered_predicate(17)
    }
    /// Computed predicate `shared.enable_simd() && x86.has_avx512dq()`.
    pub fn use_avx512dq_simd(&self) -> bool {
        self.numbered_predicate(18)
    }
    /// Computed predicate `shared.enable_simd() && x86.has_avx512f()`.
    pub fn use_avx512f_simd(&self) -> bool {
        self.numbered_predicate(19)
    }
    /// Computed predicate `shared.enable_simd() && x86.has_avx512vbmi()`.
    pub fn use_avx512vbmi_simd(&self) -> bool {
        self.numbered_predicate(20)
    }
    /// Computed predicate `shared.enable_simd() && x86.has_avx512vl()`.
    pub fn use_avx512vl_simd(&self) -> bool {
        self.numbered_predicate(21)
    }
    /// Computed predicate `shared.enable_simd() && x86.has_avx()`.
    pub fn use_avx_simd(&self) -> bool {
        self.numbered_predicate(22)
    }
    /// Computed predicate `x86.has_bmi1()`.
    pub fn use_bmi1(&self) -> bool {
        self.numbered_predicate(23)
    }
    /// Computed predicate `x86.has_avx() && x86.has_fma()`.
    pub fn use_fma(&self) -> bool {
        self.numbered_predicate(24)
    }
    /// Computed predicate `x86.has_lzcnt()`.
    pub fn use_lzcnt(&self) -> bool {
        self.numbered_predicate(25)
    }
    /// Computed predicate `x86.has_popcnt() && x86.has_sse42()`.
    pub fn use_popcnt(&self) -> bool {
        self.numbered_predicate(26)
    }
    /// Computed predicate `x86.has_sse41()`.
    pub fn use_sse41(&self) -> bool {
        self.numbered_predicate(27)
    }
    /// Computed predicate `shared.enable_simd() && x86.has_sse41()`.
    pub fn use_sse41_simd(&self) -> bool {
        self.numbered_predicate(28)
    }
    /// Computed predicate `x86.has_sse41() && x86.has_sse42()`.
    pub fn use_sse42(&self) -> bool {
        self.numbered_predicate(29)
    }
    /// Computed predicate `shared.enable_simd() && x86.has_sse41() && x86.has_sse42()`.
    pub fn use_sse42_simd(&self) -> bool {
        self.numbered_predicate(30)
    }
    /// Computed predicate `x86.has_ssse3()`.
    pub fn use_ssse3(&self) -> bool {
        self.numbered_predicate(31)
    }
    /// Computed predicate `shared.enable_simd() && x86.has_ssse3()`.
    pub fn use_ssse3_simd(&self) -> bool {
        self.numbered_predicate(32)
    }
}
static DESCRIPTORS: [detail::Descriptor; 82] = [
    detail::Descriptor {
        name: "has_sse3",
        description: "Has support for SSE3.",
        offset: 0,
        detail: detail::Detail::Bool { bit: 0 },
    },
    detail::Descriptor {
        name: "has_ssse3",
        description: "Has support for SSSE3.",
        offset: 0,
        detail: detail::Detail::Bool { bit: 1 },
    },
    detail::Descriptor {
        name: "has_sse41",
        description: "Has support for SSE4.1.",
        offset: 0,
        detail: detail::Detail::Bool { bit: 2 },
    },
    detail::Descriptor {
        name: "has_sse42",
        description: "Has support for SSE4.2.",
        offset: 0,
        detail: detail::Detail::Bool { bit: 3 },
    },
    detail::Descriptor {
        name: "has_avx",
        description: "Has support for AVX.",
        offset: 0,
        detail: detail::Detail::Bool { bit: 4 },
    },
    detail::Descriptor {
        name: "has_avx2",
        description: "Has support for AVX2.",
        offset: 0,
        detail: detail::Detail::Bool { bit: 5 },
    },
    detail::Descriptor {
        name: "has_fma",
        description: "Has support for FMA.",
        offset: 0,
        detail: detail::Detail::Bool { bit: 6 },
    },
    detail::Descriptor {
        name: "has_avx512bitalg",
        description: "Has support for AVX512BITALG.",
        offset: 0,
        detail: detail::Detail::Bool { bit: 7 },
    },
    detail::Descriptor {
        name: "has_avx512dq",
        description: "Has support for AVX512DQ.",
        offset: 1,
        detail: detail::Detail::Bool { bit: 0 },
    },
    detail::Descriptor {
        name: "has_avx512vl",
        description: "Has support for AVX512VL.",
        offset: 1,
        detail: detail::Detail::Bool { bit: 1 },
    },
    detail::Descriptor {
        name: "has_avx512vbmi",
        description: "Has support for AVX512VMBI.",
        offset: 1,
        detail: detail::Detail::Bool { bit: 2 },
    },
    detail::Descriptor {
        name: "has_avx512f",
        description: "Has support for AVX512F.",
        offset: 1,
        detail: detail::Detail::Bool { bit: 3 },
    },
    detail::Descriptor {
        name: "has_popcnt",
        description: "Has support for POPCNT.",
        offset: 1,
        detail: detail::Detail::Bool { bit: 4 },
    },
    detail::Descriptor {
        name: "has_bmi1",
        description: "Has support for BMI1.",
        offset: 1,
        detail: detail::Detail::Bool { bit: 5 },
    },
    detail::Descriptor {
        name: "has_bmi2",
        description: "Has support for BMI2.",
        offset: 1,
        detail: detail::Detail::Bool { bit: 6 },
    },
    detail::Descriptor {
        name: "has_lzcnt",
        description: "Has support for LZCNT.",
        offset: 1,
        detail: detail::Detail::Bool { bit: 7 },
    },
    detail::Descriptor {
        name: "sse3",
        description: "SSE3 and earlier.",
        offset: 0,
        detail: detail::Detail::Preset,
    },
    detail::Descriptor {
        name: "ssse3",
        description: "SSSE3 and earlier.",
        offset: 2,
        detail: detail::Detail::Preset,
    },
    detail::Descriptor {
        name: "sse41",
        description: "SSE4.1 and earlier.",
        offset: 4,
        detail: detail::Detail::Preset,
    },
    detail::Descriptor {
        name: "sse42",
        description: "SSE4.2 and earlier.",
        offset: 6,
        detail: detail::Detail::Preset,
    },
    detail::Descriptor {
        name: "baseline",
        description: "A baseline preset with no extensions enabled.",
        offset: 8,
        detail: detail::Detail::Preset,
    },
    detail::Descriptor {
        name: "nocona",
        description: "Nocona microarchitecture.",
        offset: 10,
        detail: detail::Detail::Preset,
    },
    detail::Descriptor {
        name: "core2",
        description: "Core 2 microarchitecture.",
        offset: 12,
        detail: detail::Detail::Preset,
    },
    detail::Descriptor {
        name: "penryn",
        description: "Penryn microarchitecture.",
        offset: 14,
        detail: detail::Detail::Preset,
    },
    detail::Descriptor {
        name: "atom",
        description: "Atom microarchitecture.",
        offset: 16,
        detail: detail::Detail::Preset,
    },
    detail::Descriptor {
        name: "bonnell",
        description: "Bonnell microarchitecture.",
        offset: 18,
        detail: detail::Detail::Preset,
    },
    detail::Descriptor {
        name: "silvermont",
        description: "Silvermont microarchitecture.",
        offset: 20,
        detail: detail::Detail::Preset,
    },
    detail::Descriptor {
        name: "slm",
        description: "Silvermont microarchitecture.",
        offset: 22,
        detail: detail::Detail::Preset,
    },
    detail::Descriptor {
        name: "goldmont",
        description: "Goldmont microarchitecture.",
        offset: 24,
        detail: detail::Detail::Preset,
    },
    detail::Descriptor {
        name: "goldmont-plus",
        description: "Goldmont Plus microarchitecture.",
        offset: 26,
        detail: detail::Detail::Preset,
    },
    detail::Descriptor {
        name: "tremont",
        description: "Tremont microarchitecture.",
        offset: 28,
        detail: detail::Detail::Preset,
    },
    detail::Descriptor {
        name: "alderlake",
        description: "Alderlake microarchitecture.",
        offset: 30,
        detail: detail::Detail::Preset,
    },
    detail::Descriptor {
        name: "sierraforest",
        description: "Sierra Forest microarchitecture.",
        offset: 32,
        detail: detail::Detail::Preset,
    },
    detail::Descriptor {
        name: "grandridge",
        description: "Grandridge microarchitecture.",
        offset: 34,
        detail: detail::Detail::Preset,
    },
    detail::Descriptor {
        name: "nehalem",
        description: "Nehalem microarchitecture.",
        offset: 36,
        detail: detail::Detail::Preset,
    },
    detail::Descriptor {
        name: "corei7",
        description: "Core i7 microarchitecture.",
        offset: 38,
        detail: detail::Detail::Preset,
    },
    detail::Descriptor {
        name: "westmere",
        description: "Westmere microarchitecture.",
        offset: 40,
        detail: detail::Detail::Preset,
    },
    detail::Descriptor {
        name: "sandybridge",
        description: "Sandy Bridge microarchitecture.",
        offset: 42,
        detail: detail::Detail::Preset,
    },
    detail::Descriptor {
        name: "corei7-avx",
        description: "Core i7 AVX microarchitecture.",
        offset: 44,
        detail: detail::Detail::Preset,
    },
    detail::Descriptor {
        name: "ivybridge",
        description: "Ivy Bridge microarchitecture.",
        offset: 46,
        detail: detail::Detail::Preset,
    },
    detail::Descriptor {
        name: "core-avx-i",
        description: "Intel Core CPU with 64-bit extensions.",
        offset: 48,
        detail: detail::Detail::Preset,
    },
    detail::Descriptor {
        name: "haswell",
        description: "Haswell microarchitecture.",
        offset: 50,
        detail: detail::Detail::Preset,
    },
    detail::Descriptor {
        name: "core-avx2",
        description: "Intel Core CPU with AVX2 extensions.",
        offset: 52,
        detail: detail::Detail::Preset,
    },
    detail::Descriptor {
        name: "broadwell",
        description: "Broadwell microarchitecture.",
        offset: 54,
        detail: detail::Detail::Preset,
    },
    detail::Descriptor {
        name: "skylake",
        description: "Skylake microarchitecture.",
        offset: 56,
        detail: detail::Detail::Preset,
    },
    detail::Descriptor {
        name: "knl",
        description: "Knights Landing microarchitecture.",
        offset: 58,
        detail: detail::Detail::Preset,
    },
    detail::Descriptor {
        name: "knm",
        description: "Knights Mill microarchitecture.",
        offset: 60,
        detail: detail::Detail::Preset,
    },
    detail::Descriptor {
        name: "skylake-avx512",
        description: "Skylake AVX512 microarchitecture.",
        offset: 62,
        detail: detail::Detail::Preset,
    },
    detail::Descriptor {
        name: "skx",
        description: "Skylake AVX512 microarchitecture.",
        offset: 64,
        detail: detail::Detail::Preset,
    },
    detail::Descriptor {
        name: "cascadelake",
        description: "Cascade Lake microarchitecture.",
        offset: 66,
        detail: detail::Detail::Preset,
    },
    detail::Descriptor {
        name: "cooperlake",
        description: "Cooper Lake microarchitecture.",
        offset: 68,
        detail: detail::Detail::Preset,
    },
    detail::Descriptor {
        name: "cannonlake",
        description: "Canon Lake microarchitecture.",
        offset: 70,
        detail: detail::Detail::Preset,
    },
    detail::Descriptor {
        name: "icelake-client",
        description: "Ice Lake microarchitecture.",
        offset: 72,
        detail: detail::Detail::Preset,
    },
    detail::Descriptor {
        name: "icelake",
        description: "Ice Lake microarchitecture",
        offset: 74,
        detail: detail::Detail::Preset,
    },
    detail::Descriptor {
        name: "icelake-server",
        description: "Ice Lake (server) microarchitecture.",
        offset: 76,
        detail: detail::Detail::Preset,
    },
    detail::Descriptor {
        name: "tigerlake",
        description: "Tiger Lake microarchitecture.",
        offset: 78,
        detail: detail::Detail::Preset,
    },
    detail::Descriptor {
        name: "sapphirerapids",
        description: "Saphire Rapids microarchitecture.",
        offset: 80,
        detail: detail::Detail::Preset,
    },
    detail::Descriptor {
        name: "raptorlake",
        description: "Raptor Lake microarchitecture.",
        offset: 82,
        detail: detail::Detail::Preset,
    },
    detail::Descriptor {
        name: "meteorlake",
        description: "Meteor Lake microarchitecture.",
        offset: 84,
        detail: detail::Detail::Preset,
    },
    detail::Descriptor {
        name: "graniterapids",
        description: "Granite Rapids microarchitecture.",
        offset: 86,
        detail: detail::Detail::Preset,
    },
    detail::Descriptor {
        name: "opteron",
        description: "Opteron microarchitecture.",
        offset: 88,
        detail: detail::Detail::Preset,
    },
    detail::Descriptor {
        name: "k8",
        description: "K8 Hammer microarchitecture.",
        offset: 90,
        detail: detail::Detail::Preset,
    },
    detail::Descriptor {
        name: "athlon64",
        description: "Athlon64 microarchitecture.",
        offset: 92,
        detail: detail::Detail::Preset,
    },
    detail::Descriptor {
        name: "athlon-fx",
        description: "Athlon FX microarchitecture.",
        offset: 94,
        detail: detail::Detail::Preset,
    },
    detail::Descriptor {
        name: "opteron-sse3",
        description: "Opteron microarchitecture with support for SSE3 instructions.",
        offset: 96,
        detail: detail::Detail::Preset,
    },
    detail::Descriptor {
        name: "k8-sse3",
        description: "K8 Hammer microarchitecture with support for SSE3 instructions.",
        offset: 98,
        detail: detail::Detail::Preset,
    },
    detail::Descriptor {
        name: "athlon64-sse3",
        description: "Athlon 64 microarchitecture with support for SSE3 instructions.",
        offset: 100,
        detail: detail::Detail::Preset,
    },
    detail::Descriptor {
        name: "barcelona",
        description: "Barcelona microarchitecture.",
        offset: 102,
        detail: detail::Detail::Preset,
    },
    detail::Descriptor {
        name: "amdfam10",
        description: "AMD Family 10h microarchitecture",
        offset: 104,
        detail: detail::Detail::Preset,
    },
    detail::Descriptor {
        name: "btver1",
        description: "Bobcat microarchitecture.",
        offset: 106,
        detail: detail::Detail::Preset,
    },
    detail::Descriptor {
        name: "btver2",
        description: "Jaguar microarchitecture.",
        offset: 108,
        detail: detail::Detail::Preset,
    },
    detail::Descriptor {
        name: "bdver1",
        description: "Bulldozer microarchitecture",
        offset: 110,
        detail: detail::Detail::Preset,
    },
    detail::Descriptor {
        name: "bdver2",
        description: "Piledriver microarchitecture.",
        offset: 112,
        detail: detail::Detail::Preset,
    },
    detail::Descriptor {
        name: "bdver3",
        description: "Steamroller microarchitecture.",
        offset: 114,
        detail: detail::Detail::Preset,
    },
    detail::Descriptor {
        name: "bdver4",
        description: "Excavator microarchitecture.",
        offset: 116,
        detail: detail::Detail::Preset,
    },
    detail::Descriptor {
        name: "znver1",
        description: "Zen (first generation) microarchitecture.",
        offset: 118,
        detail: detail::Detail::Preset,
    },
    detail::Descriptor {
        name: "znver2",
        description: "Zen (second generation) microarchitecture.",
        offset: 120,
        detail: detail::Detail::Preset,
    },
    detail::Descriptor {
        name: "znver3",
        description: "Zen (third generation) microarchitecture.",
        offset: 122,
        detail: detail::Detail::Preset,
    },
    detail::Descriptor {
        name: "x86-64",
        description: "Generic x86-64 microarchitecture.",
        offset: 124,
        detail: detail::Detail::Preset,
    },
    detail::Descriptor {
        name: "x86-64-v2",
        description: "Generic x86-64 (V2) microarchitecture.",
        offset: 126,
        detail: detail::Detail::Preset,
    },
    detail::Descriptor {
        name: "x84_64_v3",
        description: "Generic x86_64 (V3) microarchitecture.",
        offset: 128,
        detail: detail::Detail::Preset,
    },
    detail::Descriptor {
        name: "x86_64_v4",
        description: "Generic x86_64 (V4) microarchitecture.",
        offset: 130,
        detail: detail::Detail::Preset,
    },
];
static ENUMERATORS: [&str; 0] = [
];
static HASH_TABLE: [u16; 128] = [
    0xffff,
    0xffff,
    77,
    76,
    75,
    0xffff,
    0xffff,
    0xffff,
    23,
    0xffff,
    66,
    79,
    22,
    50,
    59,
    14,
    13,
    29,
    1,
    41,
    70,
    67,
    4,
    35,
    0xffff,
    65,
    5,
    44,
    21,
    64,
    15,
    6,
    47,
    49,
    24,
    62,
    0xffff,
    11,
    43,
    38,
    52,
    0xffff,
    0xffff,
    69,
    0xffff,
    3,
    31,
    0xffff,
    2,
    0xffff,
    0xffff,
    58,
    0xffff,
    0xffff,
    10,
    12,
    0xffff,
    0xffff,
    0xffff,
    0xffff,
    0xffff,
    0xffff,
    0xffff,
    0xffff,
    30,
    78,
    73,
    0,
    39,
    28,
    46,
    45,
    8,
    54,
    71,
    9,
    74,
    72,
    0xffff,
    0xffff,
    0xffff,
    61,
    80,
    33,
    7,
    0xffff,
    18,
    19,
    48,
    16,
    53,
    60,
    0xffff,
    0xffff,
    20,
    0xffff,
    63,
    68,
    56,
    0xffff,
    0xffff,
    81,
    0xffff,
    26,
    27,
    0xffff,
    34,
    0xffff,
    0xffff,
    36,
    0xffff,
    0xffff,
    40,
    42,
    0xffff,
    32,
    0xffff,
    0xffff,
    0xffff,
    57,
    51,
    0xffff,
    0xffff,
    17,
    55,
    0xffff,
    25,
    37,
];
static PRESETS: [(u8, u8); 132] = [
    // sse3: has_sse3
    (0b00000001, 0b00000001),
    (0b00000000, 0b00000000),
    // ssse3: has_sse3, has_ssse3
    (0b00000011, 0b00000011),
    (0b00000000, 0b00000000),
    // sse41: has_sse3, has_ssse3, has_sse41
    (0b00000111, 0b00000111),
    (0b00000000, 0b00000000),
    // sse42: has_sse3, has_ssse3, has_sse41, has_sse42
    (0b00001111, 0b00001111),
    (0b00000000, 0b00000000),
    // baseline: 
    (0b00000000, 0b00000000),
    (0b00000000, 0b00000000),
    // nocona: has_sse3
    (0b00000001, 0b00000001),
    (0b00000000, 0b00000000),
    // core2: has_sse3
    (0b00000001, 0b00000001),
    (0b00000000, 0b00000000),
    // penryn: has_sse3, has_ssse3, has_sse41
    (0b00000111, 0b00000111),
    (0b00000000, 0b00000000),
    // atom: has_sse3, has_ssse3
    (0b00000011, 0b00000011),
    (0b00000000, 0b00000000),
    // bonnell: has_sse3, has_ssse3
    (0b00000011, 0b00000011),
    (0b00000000, 0b00000000),
    // silvermont: has_sse3, has_ssse3, has_sse3, has_ssse3, has_sse41, has_sse42, has_popcnt
    (0b00001111, 0b00001111),
    (0b00010000, 0b00010000),
    // slm: has_sse3, has_ssse3, has_sse3, has_ssse3, has_sse41, has_sse42, has_popcnt
    (0b00001111, 0b00001111),
    (0b00010000, 0b00010000),
    // goldmont: has_sse3, has_ssse3, has_sse3, has_ssse3, has_sse41, has_sse42, has_popcnt
    (0b00001111, 0b00001111),
    (0b00010000, 0b00010000),
    // goldmont-plus: has_sse3, has_ssse3, has_sse3, has_ssse3, has_sse41, has_sse42, has_popcnt
    (0b00001111, 0b00001111),
    (0b00010000, 0b00010000),
    // tremont: has_sse3, has_ssse3, has_sse3, has_ssse3, has_sse41, has_sse42, has_popcnt
    (0b00001111, 0b00001111),
    (0b00010000, 0b00010000),
    // alderlake: has_sse3, has_ssse3, has_sse3, has_ssse3, has_sse41, has_sse42, has_popcnt, has_bmi1, has_bmi2, has_lzcnt, has_fma
    (0b01001111, 0b01001111),
    (0b11110000, 0b11110000),
    // sierraforest: has_sse3, has_ssse3, has_sse3, has_ssse3, has_sse41, has_sse42, has_popcnt, has_bmi1, has_bmi2, has_lzcnt, has_fma
    (0b01001111, 0b01001111),
    (0b11110000, 0b11110000),
    // grandridge: has_sse3, has_ssse3, has_sse3, has_ssse3, has_sse41, has_sse42, has_popcnt, has_bmi1, has_bmi2, has_lzcnt, has_fma
    (0b01001111, 0b01001111),
    (0b11110000, 0b11110000),
    // nehalem: has_sse3, has_ssse3, has_sse41, has_sse42, has_popcnt
    (0b00001111, 0b00001111),
    (0b00010000, 0b00010000),
    // corei7: has_sse3, has_ssse3, has_sse41, has_sse42, has_popcnt
    (0b00001111, 0b00001111),
    (0b00010000, 0b00010000),
    // westmere: has_sse3, has_ssse3, has_sse41, has_sse42, has_popcnt
    (0b00001111, 0b00001111),
    (0b00010000, 0b00010000),
    // sandybridge: has_sse3, has_ssse3, has_sse41, has_sse42, has_popcnt, has_avx
    (0b00011111, 0b00011111),
    (0b00010000, 0b00010000),
    // corei7-avx: has_sse3, has_ssse3, has_sse41, has_sse42, has_popcnt, has_avx
    (0b00011111, 0b00011111),
    (0b00010000, 0b00010000),
    // ivybridge: has_sse3, has_ssse3, has_sse41, has_sse42, has_popcnt, has_avx
    (0b00011111, 0b00011111),
    (0b00010000, 0b00010000),
    // core-avx-i: has_sse3, has_ssse3, has_sse41, has_sse42, has_popcnt, has_avx
    (0b00011111, 0b00011111),
    (0b00010000, 0b00010000),
    // haswell: has_sse3, has_ssse3, has_sse41, has_sse42, has_popcnt, has_avx, has_avx2, has_bmi1, has_bmi2, has_fma, has_lzcnt
    (0b01111111, 0b01111111),
    (0b11110000, 0b11110000),
    // core-avx2: has_sse3, has_ssse3, has_sse41, has_sse42, has_popcnt, has_avx, has_avx2, has_bmi1, has_bmi2, has_fma, has_lzcnt
    (0b01111111, 0b01111111),
    (0b11110000, 0b11110000),
    // broadwell: has_sse3, has_ssse3, has_sse41, has_sse42, has_popcnt, has_avx, has_avx2, has_bmi1, has_bmi2, has_fma, has_lzcnt
    (0b01111111, 0b01111111),
    (0b11110000, 0b11110000),
    // skylake: has_sse3, has_ssse3, has_sse41, has_sse42, has_popcnt, has_avx, has_avx2, has_bmi1, has_bmi2, has_fma, has_lzcnt
    (0b01111111, 0b01111111),
    (0b11110000, 0b11110000),
    // knl: has_popcnt, has_avx512f, has_fma, has_bmi1, has_bmi2, has_lzcnt
    (0b01000000, 0b01000000),
    (0b11111000, 0b11111000),
    // knm: has_popcnt, has_avx512f, has_fma, has_bmi1, has_bmi2, has_lzcnt
    (0b01000000, 0b01000000),
    (0b11111000, 0b11111000),
    // skylake-avx512: has_sse3, has_ssse3, has_sse41, has_sse42, has_popcnt, has_avx, has_avx2, has_bmi1, has_bmi2, has_fma, has_lzcnt, has_avx512f, has_avx512dq, has_avx512vl
    (0b01111111, 0b01111111),
    (0b11111011, 0b11111011),
    // skx: has_sse3, has_ssse3, has_sse41, has_sse42, has_popcnt, has_avx, has_avx2, has_bmi1, has_bmi2, has_fma, has_lzcnt, has_avx512f, has_avx512dq, has_avx512vl
    (0b01111111, 0b01111111),
    (0b11111011, 0b11111011),
    // cascadelake: has_sse3, has_ssse3, has_sse41, has_sse42, has_popcnt, has_avx, has_avx2, has_bmi1, has_bmi2, has_fma, has_lzcnt, has_avx512f, has_avx512dq, has_avx512vl
    (0b01111111, 0b01111111),
    (0b11111011, 0b11111011),
    // cooperlake: has_sse3, has_ssse3, has_sse41, has_sse42, has_popcnt, has_avx, has_avx2, has_bmi1, has_bmi2, has_fma, has_lzcnt, has_avx512f, has_avx512dq, has_avx512vl
    (0b01111111, 0b01111111),
    (0b11111011, 0b11111011),
    // cannonlake: has_sse3, has_ssse3, has_sse41, has_sse42, has_popcnt, has_avx, has_avx2, has_bmi1, has_bmi2, has_fma, has_lzcnt, has_avx512f, has_avx512dq, has_avx512vl, has_avx512vbmi
    (0b01111111, 0b01111111),
    (0b11111111, 0b11111111),
    // icelake-client: has_sse3, has_ssse3, has_sse41, has_sse42, has_popcnt, has_avx, has_avx2, has_bmi1, has_bmi2, has_fma, has_lzcnt, has_avx512f, has_avx512dq, has_avx512vl, has_avx512vbmi, has_avx512bitalg
    (0b11111111, 0b11111111),
    (0b11111111, 0b11111111),
    // icelake: has_sse3, has_ssse3, has_sse41, has_sse42, has_popcnt, has_avx, has_avx2, has_bmi1, has_bmi2, has_fma, has_lzcnt, has_avx512f, has_avx512dq, has_avx512vl, has_avx512vbmi, has_avx512bitalg
    (0b11111111, 0b11111111),
    (0b11111111, 0b11111111),
    // icelake-server: has_sse3, has_ssse3, has_sse41, has_sse42, has_popcnt, has_avx, has_avx2, has_bmi1, has_bmi2, has_fma, has_lzcnt, has_avx512f, has_avx512dq, has_avx512vl, has_avx512vbmi, has_avx512bitalg
    (0b11111111, 0b11111111),
    (0b11111111, 0b11111111),
    // tigerlake: has_sse3, has_ssse3, has_sse41, has_sse42, has_popcnt, has_avx, has_avx2, has_bmi1, has_bmi2, has_fma, has_lzcnt, has_avx512f, has_avx512dq, has_avx512vl, has_avx512vbmi, has_avx512bitalg
    (0b11111111, 0b11111111),
    (0b11111111, 0b11111111),
    // sapphirerapids: has_sse3, has_ssse3, has_sse41, has_sse42, has_popcnt, has_avx, has_avx2, has_bmi1, has_bmi2, has_fma, has_lzcnt, has_avx512f, has_avx512dq, has_avx512vl, has_avx512vbmi, has_avx512bitalg
    (0b11111111, 0b11111111),
    (0b11111111, 0b11111111),
    // raptorlake: has_sse3, has_ssse3, has_sse3, has_ssse3, has_sse41, has_sse42, has_popcnt, has_bmi1, has_bmi2, has_lzcnt, has_fma
    (0b01001111, 0b01001111),
    (0b11110000, 0b11110000),
    // meteorlake: has_sse3, has_ssse3, has_sse3, has_ssse3, has_sse41, has_sse42, has_popcnt, has_bmi1, has_bmi2, has_lzcnt, has_fma
    (0b01001111, 0b01001111),
    (0b11110000, 0b11110000),
    // graniterapids: has_sse3, has_ssse3, has_sse41, has_sse42, has_popcnt, has_avx, has_avx2, has_bmi1, has_bmi2, has_fma, has_lzcnt, has_avx512f, has_avx512dq, has_avx512vl, has_avx512vbmi, has_avx512bitalg
    (0b11111111, 0b11111111),
    (0b11111111, 0b11111111),
    // opteron: 
    (0b00000000, 0b00000000),
    (0b00000000, 0b00000000),
    // k8: 
    (0b00000000, 0b00000000),
    (0b00000000, 0b00000000),
    // athlon64: 
    (0b00000000, 0b00000000),
    (0b00000000, 0b00000000),
    // athlon-fx: 
    (0b00000000, 0b00000000),
    (0b00000000, 0b00000000),
    // opteron-sse3: has_sse3
    (0b00000001, 0b00000001),
    (0b00000000, 0b00000000),
    // k8-sse3: has_sse3
    (0b00000001, 0b00000001),
    (0b00000000, 0b00000000),
    // athlon64-sse3: has_sse3
    (0b00000001, 0b00000001),
    (0b00000000, 0b00000000),
    // barcelona: has_popcnt, has_lzcnt
    (0b00000000, 0b00000000),
    (0b10010000, 0b10010000),
    // amdfam10: has_popcnt, has_lzcnt
    (0b00000000, 0b00000000),
    (0b10010000, 0b10010000),
    // btver1: has_sse3, has_ssse3, has_lzcnt, has_popcnt
    (0b00000011, 0b00000011),
    (0b10010000, 0b10010000),
    // btver2: has_sse3, has_ssse3, has_lzcnt, has_popcnt, has_avx, has_bmi1
    (0b00010011, 0b00010011),
    (0b10110000, 0b10110000),
    // bdver1: has_lzcnt, has_popcnt, has_sse3, has_ssse3
    (0b00000011, 0b00000011),
    (0b10010000, 0b10010000),
    // bdver2: has_lzcnt, has_popcnt, has_sse3, has_ssse3, has_bmi1
    (0b00000011, 0b00000011),
    (0b10110000, 0b10110000),
    // bdver3: has_lzcnt, has_popcnt, has_sse3, has_ssse3, has_bmi1
    (0b00000011, 0b00000011),
    (0b10110000, 0b10110000),
    // bdver4: has_lzcnt, has_popcnt, has_sse3, has_ssse3, has_bmi1, has_avx2, has_bmi2
    (0b00100011, 0b00100011),
    (0b11110000, 0b11110000),
    // znver1: has_sse3, has_ssse3, has_sse41, has_sse42, has_popcnt, has_bmi1, has_bmi2, has_lzcnt, has_fma
    (0b01001111, 0b01001111),
    (0b11110000, 0b11110000),
    // znver2: has_sse3, has_ssse3, has_sse41, has_sse42, has_popcnt, has_bmi1, has_bmi2, has_lzcnt, has_fma
    (0b01001111, 0b01001111),
    (0b11110000, 0b11110000),
    // znver3: has_sse3, has_ssse3, has_sse41, has_sse42, has_popcnt, has_bmi1, has_bmi2, has_lzcnt, has_fma
    (0b01001111, 0b01001111),
    (0b11110000, 0b11110000),
    // x86-64: 
    (0b00000000, 0b00000000),
    (0b00000000, 0b00000000),
    // x86-64-v2: has_sse3, has_ssse3, has_sse41, has_sse42, has_popcnt
    (0b00001111, 0b00001111),
    (0b00010000, 0b00010000),
    // x84_64_v3: has_sse3, has_ssse3, has_sse41, has_sse42, has_popcnt, has_bmi1, has_bmi2, has_fma, has_lzcnt, has_avx2
    (0b01101111, 0b01101111),
    (0b11110000, 0b11110000),
    // x86_64_v4: has_sse3, has_ssse3, has_sse41, has_sse42, has_popcnt, has_bmi1, has_bmi2, has_fma, has_lzcnt, has_avx2, has_avx512dq, has_avx512vl
    (0b01101111, 0b01101111),
    (0b11110011, 0b11110011),
];
static TEMPLATE: detail::Template = detail::Template {
    name: "x86",
    descriptors: &DESCRIPTORS,
    enumerators: &ENUMERATORS,
    hash_table: &HASH_TABLE,
    defaults: &[0x0f, 0x00],
    presets: &PRESETS,
};
/// Create a `settings::Builder` for the x86 settings group.
pub fn builder() -> Builder {
    Builder::new(&TEMPLATE)
}
impl fmt::Display for Flags {
    fn fmt(&self, f: &mut fmt::Formatter) -> fmt::Result {
        writeln!(f, "[x86]")?;
        for d in &DESCRIPTORS {
            if !d.detail.is_preset() {
                write!(f, "{} = ", d.name)?;
                TEMPLATE.format_toml_value(d.detail, self.bytes[d.offset as usize], f)?;
                writeln!(f)?;
            }
        }
        Ok(())
    }
}
