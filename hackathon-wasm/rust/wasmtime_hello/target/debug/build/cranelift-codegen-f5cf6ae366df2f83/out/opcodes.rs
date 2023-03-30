/// An instruction format
///
/// Every opcode has a corresponding instruction format
/// which is represented by both the `InstructionFormat`
/// and the `InstructionData` enums.
#[derive(Copy, Clone, PartialEq, Eq, Debug)]
pub enum InstructionFormat {
    /// AtomicCas(imms=(flags: ir::MemFlags), vals=3, blocks=0)
    AtomicCas,
    /// AtomicRmw(imms=(flags: ir::MemFlags, op: ir::AtomicRmwOp), vals=2, blocks=0)
    AtomicRmw,
    /// Binary(imms=(), vals=2, blocks=0)
    Binary,
    /// BinaryImm64(imms=(imm: ir::immediates::Imm64), vals=1, blocks=0)
    BinaryImm64,
    /// BinaryImm8(imms=(imm: ir::immediates::Uimm8), vals=1, blocks=0)
    BinaryImm8,
    /// BranchTable(imms=(table: ir::JumpTable), vals=1, blocks=0)
    BranchTable,
    /// Brif(imms=(), vals=1, blocks=2)
    Brif,
    /// Call(imms=(func_ref: ir::FuncRef), vals=0, blocks=0)
    Call,
    /// CallIndirect(imms=(sig_ref: ir::SigRef), vals=1, blocks=0)
    CallIndirect,
    /// CondTrap(imms=(code: ir::TrapCode), vals=1, blocks=0)
    CondTrap,
    /// DynamicStackLoad(imms=(dynamic_stack_slot: ir::DynamicStackSlot), vals=0, blocks=0)
    DynamicStackLoad,
    /// DynamicStackStore(imms=(dynamic_stack_slot: ir::DynamicStackSlot), vals=1, blocks=0)
    DynamicStackStore,
    /// FloatCompare(imms=(cond: ir::condcodes::FloatCC), vals=2, blocks=0)
    FloatCompare,
    /// FuncAddr(imms=(func_ref: ir::FuncRef), vals=0, blocks=0)
    FuncAddr,
    /// IntAddTrap(imms=(code: ir::TrapCode), vals=2, blocks=0)
    IntAddTrap,
    /// IntCompare(imms=(cond: ir::condcodes::IntCC), vals=2, blocks=0)
    IntCompare,
    /// IntCompareImm(imms=(cond: ir::condcodes::IntCC, imm: ir::immediates::Imm64), vals=1, blocks=0)
    IntCompareImm,
    /// Jump(imms=(), vals=0, blocks=1)
    Jump,
    /// Load(imms=(flags: ir::MemFlags, offset: ir::immediates::Offset32), vals=1, blocks=0)
    Load,
    /// LoadNoOffset(imms=(flags: ir::MemFlags), vals=1, blocks=0)
    LoadNoOffset,
    /// MultiAry(imms=(), vals=0, blocks=0)
    MultiAry,
    /// NullAry(imms=(), vals=0, blocks=0)
    NullAry,
    /// Shuffle(imms=(imm: ir::Immediate), vals=2, blocks=0)
    Shuffle,
    /// StackLoad(imms=(stack_slot: ir::StackSlot, offset: ir::immediates::Offset32), vals=0, blocks=0)
    StackLoad,
    /// StackStore(imms=(stack_slot: ir::StackSlot, offset: ir::immediates::Offset32), vals=1, blocks=0)
    StackStore,
    /// Store(imms=(flags: ir::MemFlags, offset: ir::immediates::Offset32), vals=2, blocks=0)
    Store,
    /// StoreNoOffset(imms=(flags: ir::MemFlags), vals=2, blocks=0)
    StoreNoOffset,
    /// TableAddr(imms=(table: ir::Table, offset: ir::immediates::Offset32), vals=1, blocks=0)
    TableAddr,
    /// Ternary(imms=(), vals=3, blocks=0)
    Ternary,
    /// TernaryImm8(imms=(imm: ir::immediates::Uimm8), vals=2, blocks=0)
    TernaryImm8,
    /// Trap(imms=(code: ir::TrapCode), vals=0, blocks=0)
    Trap,
    /// Unary(imms=(), vals=1, blocks=0)
    Unary,
    /// UnaryConst(imms=(constant_handle: ir::Constant), vals=0, blocks=0)
    UnaryConst,
    /// UnaryGlobalValue(imms=(global_value: ir::GlobalValue), vals=0, blocks=0)
    UnaryGlobalValue,
    /// UnaryIeee32(imms=(imm: ir::immediates::Ieee32), vals=0, blocks=0)
    UnaryIeee32,
    /// UnaryIeee64(imms=(imm: ir::immediates::Ieee64), vals=0, blocks=0)
    UnaryIeee64,
    /// UnaryImm(imms=(imm: ir::immediates::Imm64), vals=0, blocks=0)
    UnaryImm,
}

impl<'a> From<&'a InstructionData> for InstructionFormat {
    fn from(inst: &'a InstructionData) -> Self {
        match *inst {
            InstructionData::AtomicCas { .. } => {
                Self::AtomicCas
            }
            InstructionData::AtomicRmw { .. } => {
                Self::AtomicRmw
            }
            InstructionData::Binary { .. } => {
                Self::Binary
            }
            InstructionData::BinaryImm64 { .. } => {
                Self::BinaryImm64
            }
            InstructionData::BinaryImm8 { .. } => {
                Self::BinaryImm8
            }
            InstructionData::BranchTable { .. } => {
                Self::BranchTable
            }
            InstructionData::Brif { .. } => {
                Self::Brif
            }
            InstructionData::Call { .. } => {
                Self::Call
            }
            InstructionData::CallIndirect { .. } => {
                Self::CallIndirect
            }
            InstructionData::CondTrap { .. } => {
                Self::CondTrap
            }
            InstructionData::DynamicStackLoad { .. } => {
                Self::DynamicStackLoad
            }
            InstructionData::DynamicStackStore { .. } => {
                Self::DynamicStackStore
            }
            InstructionData::FloatCompare { .. } => {
                Self::FloatCompare
            }
            InstructionData::FuncAddr { .. } => {
                Self::FuncAddr
            }
            InstructionData::IntAddTrap { .. } => {
                Self::IntAddTrap
            }
            InstructionData::IntCompare { .. } => {
                Self::IntCompare
            }
            InstructionData::IntCompareImm { .. } => {
                Self::IntCompareImm
            }
            InstructionData::Jump { .. } => {
                Self::Jump
            }
            InstructionData::Load { .. } => {
                Self::Load
            }
            InstructionData::LoadNoOffset { .. } => {
                Self::LoadNoOffset
            }
            InstructionData::MultiAry { .. } => {
                Self::MultiAry
            }
            InstructionData::NullAry { .. } => {
                Self::NullAry
            }
            InstructionData::Shuffle { .. } => {
                Self::Shuffle
            }
            InstructionData::StackLoad { .. } => {
                Self::StackLoad
            }
            InstructionData::StackStore { .. } => {
                Self::StackStore
            }
            InstructionData::Store { .. } => {
                Self::Store
            }
            InstructionData::StoreNoOffset { .. } => {
                Self::StoreNoOffset
            }
            InstructionData::TableAddr { .. } => {
                Self::TableAddr
            }
            InstructionData::Ternary { .. } => {
                Self::Ternary
            }
            InstructionData::TernaryImm8 { .. } => {
                Self::TernaryImm8
            }
            InstructionData::Trap { .. } => {
                Self::Trap
            }
            InstructionData::Unary { .. } => {
                Self::Unary
            }
            InstructionData::UnaryConst { .. } => {
                Self::UnaryConst
            }
            InstructionData::UnaryGlobalValue { .. } => {
                Self::UnaryGlobalValue
            }
            InstructionData::UnaryIeee32 { .. } => {
                Self::UnaryIeee32
            }
            InstructionData::UnaryIeee64 { .. } => {
                Self::UnaryIeee64
            }
            InstructionData::UnaryImm { .. } => {
                Self::UnaryImm
            }
        }
    }
}

#[derive(Copy, Clone, Debug, PartialEq, Eq, Hash)]
#[cfg_attr(feature = "enable-serde", derive(Serialize, Deserialize))]
#[allow(missing_docs)]
pub enum InstructionData {
    AtomicCas {
        opcode: Opcode,
        args: [Value; 3],
        flags: ir::MemFlags,
    },
    AtomicRmw {
        opcode: Opcode,
        args: [Value; 2],
        flags: ir::MemFlags,
        op: ir::AtomicRmwOp,
    },
    Binary {
        opcode: Opcode,
        args: [Value; 2],
    },
    BinaryImm64 {
        opcode: Opcode,
        arg: Value,
        imm: ir::immediates::Imm64,
    },
    BinaryImm8 {
        opcode: Opcode,
        arg: Value,
        imm: ir::immediates::Uimm8,
    },
    BranchTable {
        opcode: Opcode,
        arg: Value,
        table: ir::JumpTable,
    },
    Brif {
        opcode: Opcode,
        arg: Value,
        blocks: [ir::BlockCall; 2],
    },
    Call {
        opcode: Opcode,
        args: ValueList,
        func_ref: ir::FuncRef,
    },
    CallIndirect {
        opcode: Opcode,
        args: ValueList,
        sig_ref: ir::SigRef,
    },
    CondTrap {
        opcode: Opcode,
        arg: Value,
        code: ir::TrapCode,
    },
    DynamicStackLoad {
        opcode: Opcode,
        dynamic_stack_slot: ir::DynamicStackSlot,
    },
    DynamicStackStore {
        opcode: Opcode,
        arg: Value,
        dynamic_stack_slot: ir::DynamicStackSlot,
    },
    FloatCompare {
        opcode: Opcode,
        args: [Value; 2],
        cond: ir::condcodes::FloatCC,
    },
    FuncAddr {
        opcode: Opcode,
        func_ref: ir::FuncRef,
    },
    IntAddTrap {
        opcode: Opcode,
        args: [Value; 2],
        code: ir::TrapCode,
    },
    IntCompare {
        opcode: Opcode,
        args: [Value; 2],
        cond: ir::condcodes::IntCC,
    },
    IntCompareImm {
        opcode: Opcode,
        arg: Value,
        cond: ir::condcodes::IntCC,
        imm: ir::immediates::Imm64,
    },
    Jump {
        opcode: Opcode,
        destination: ir::BlockCall,
    },
    Load {
        opcode: Opcode,
        arg: Value,
        flags: ir::MemFlags,
        offset: ir::immediates::Offset32,
    },
    LoadNoOffset {
        opcode: Opcode,
        arg: Value,
        flags: ir::MemFlags,
    },
    MultiAry {
        opcode: Opcode,
        args: ValueList,
    },
    NullAry {
        opcode: Opcode,
    },
    Shuffle {
        opcode: Opcode,
        args: [Value; 2],
        imm: ir::Immediate,
    },
    StackLoad {
        opcode: Opcode,
        stack_slot: ir::StackSlot,
        offset: ir::immediates::Offset32,
    },
    StackStore {
        opcode: Opcode,
        arg: Value,
        stack_slot: ir::StackSlot,
        offset: ir::immediates::Offset32,
    },
    Store {
        opcode: Opcode,
        args: [Value; 2],
        flags: ir::MemFlags,
        offset: ir::immediates::Offset32,
    },
    StoreNoOffset {
        opcode: Opcode,
        args: [Value; 2],
        flags: ir::MemFlags,
    },
    TableAddr {
        opcode: Opcode,
        arg: Value,
        table: ir::Table,
        offset: ir::immediates::Offset32,
    },
    Ternary {
        opcode: Opcode,
        args: [Value; 3],
    },
    TernaryImm8 {
        opcode: Opcode,
        args: [Value; 2],
        imm: ir::immediates::Uimm8,
    },
    Trap {
        opcode: Opcode,
        code: ir::TrapCode,
    },
    Unary {
        opcode: Opcode,
        arg: Value,
    },
    UnaryConst {
        opcode: Opcode,
        constant_handle: ir::Constant,
    },
    UnaryGlobalValue {
        opcode: Opcode,
        global_value: ir::GlobalValue,
    },
    UnaryIeee32 {
        opcode: Opcode,
        imm: ir::immediates::Ieee32,
    },
    UnaryIeee64 {
        opcode: Opcode,
        imm: ir::immediates::Ieee64,
    },
    UnaryImm {
        opcode: Opcode,
        imm: ir::immediates::Imm64,
    },
}

impl InstructionData {
    /// Get the opcode of this instruction.
    pub fn opcode(&self) -> Opcode {
        match *self {
            Self::AtomicCas { opcode, .. } |
            Self::AtomicRmw { opcode, .. } |
            Self::Binary { opcode, .. } |
            Self::BinaryImm64 { opcode, .. } |
            Self::BinaryImm8 { opcode, .. } |
            Self::BranchTable { opcode, .. } |
            Self::Brif { opcode, .. } |
            Self::Call { opcode, .. } |
            Self::CallIndirect { opcode, .. } |
            Self::CondTrap { opcode, .. } |
            Self::DynamicStackLoad { opcode, .. } |
            Self::DynamicStackStore { opcode, .. } |
            Self::FloatCompare { opcode, .. } |
            Self::FuncAddr { opcode, .. } |
            Self::IntAddTrap { opcode, .. } |
            Self::IntCompare { opcode, .. } |
            Self::IntCompareImm { opcode, .. } |
            Self::Jump { opcode, .. } |
            Self::Load { opcode, .. } |
            Self::LoadNoOffset { opcode, .. } |
            Self::MultiAry { opcode, .. } |
            Self::NullAry { opcode, .. } |
            Self::Shuffle { opcode, .. } |
            Self::StackLoad { opcode, .. } |
            Self::StackStore { opcode, .. } |
            Self::Store { opcode, .. } |
            Self::StoreNoOffset { opcode, .. } |
            Self::TableAddr { opcode, .. } |
            Self::Ternary { opcode, .. } |
            Self::TernaryImm8 { opcode, .. } |
            Self::Trap { opcode, .. } |
            Self::Unary { opcode, .. } |
            Self::UnaryConst { opcode, .. } |
            Self::UnaryGlobalValue { opcode, .. } |
            Self::UnaryIeee32 { opcode, .. } |
            Self::UnaryIeee64 { opcode, .. } |
            Self::UnaryImm { opcode, .. } => {
                opcode
            }
        }
    }

    /// Get the controlling type variable operand.
    pub fn typevar_operand(&self, pool: &ir::ValueListPool) -> Option<Value> {
        match *self {
            Self::Call { .. } |
            Self::DynamicStackLoad { .. } |
            Self::FuncAddr { .. } |
            Self::Jump { .. } |
            Self::MultiAry { .. } |
            Self::NullAry { .. } |
            Self::StackLoad { .. } |
            Self::Trap { .. } |
            Self::UnaryConst { .. } |
            Self::UnaryGlobalValue { .. } |
            Self::UnaryIeee32 { .. } |
            Self::UnaryIeee64 { .. } |
            Self::UnaryImm { .. } => {
                None
            }
            Self::BinaryImm64 { arg, .. } |
            Self::BinaryImm8 { arg, .. } |
            Self::BranchTable { arg, .. } |
            Self::Brif { arg, .. } |
            Self::CondTrap { arg, .. } |
            Self::DynamicStackStore { arg, .. } |
            Self::IntCompareImm { arg, .. } |
            Self::Load { arg, .. } |
            Self::LoadNoOffset { arg, .. } |
            Self::StackStore { arg, .. } |
            Self::TableAddr { arg, .. } |
            Self::Unary { arg, .. } => {
                Some(arg)
            }
            Self::AtomicRmw { args: ref args_arity2, .. } |
            Self::Binary { args: ref args_arity2, .. } |
            Self::FloatCompare { args: ref args_arity2, .. } |
            Self::IntAddTrap { args: ref args_arity2, .. } |
            Self::IntCompare { args: ref args_arity2, .. } |
            Self::Shuffle { args: ref args_arity2, .. } |
            Self::Store { args: ref args_arity2, .. } |
            Self::StoreNoOffset { args: ref args_arity2, .. } |
            Self::TernaryImm8 { args: ref args_arity2, .. } => {
                Some(args_arity2[0])
            }
            Self::Ternary { args: ref args_arity3, .. } => {
                Some(args_arity3[1])
            }
            Self::AtomicCas { args: ref args_arity3, .. } => {
                Some(args_arity3[2])
            }
            Self::CallIndirect { ref args, .. } => {
                args.get(0, pool)
            }
        }
    }

    /// Get the value arguments to this instruction.
    pub fn arguments<'a>(&'a self, pool: &'a ir::ValueListPool) -> &[Value] {
        match *self {
            Self::DynamicStackLoad { .. } |
            Self::FuncAddr { .. } |
            Self::Jump { .. } |
            Self::NullAry { .. } |
            Self::StackLoad { .. } |
            Self::Trap { .. } |
            Self::UnaryConst { .. } |
            Self::UnaryGlobalValue { .. } |
            Self::UnaryIeee32 { .. } |
            Self::UnaryIeee64 { .. } |
            Self::UnaryImm { .. } => {
                &[]
            }
            Self::AtomicRmw { args: ref args_arity2, .. } |
            Self::Binary { args: ref args_arity2, .. } |
            Self::FloatCompare { args: ref args_arity2, .. } |
            Self::IntAddTrap { args: ref args_arity2, .. } |
            Self::IntCompare { args: ref args_arity2, .. } |
            Self::Shuffle { args: ref args_arity2, .. } |
            Self::Store { args: ref args_arity2, .. } |
            Self::StoreNoOffset { args: ref args_arity2, .. } |
            Self::TernaryImm8 { args: ref args_arity2, .. } => {
                args_arity2
            }
            Self::AtomicCas { args: ref args_arity3, .. } |
            Self::Ternary { args: ref args_arity3, .. } => {
                args_arity3
            }
            Self::BinaryImm64 { ref arg, .. } |
            Self::BinaryImm8 { ref arg, .. } |
            Self::BranchTable { ref arg, .. } |
            Self::Brif { ref arg, .. } |
            Self::CondTrap { ref arg, .. } |
            Self::DynamicStackStore { ref arg, .. } |
            Self::IntCompareImm { ref arg, .. } |
            Self::Load { ref arg, .. } |
            Self::LoadNoOffset { ref arg, .. } |
            Self::StackStore { ref arg, .. } |
            Self::TableAddr { ref arg, .. } |
            Self::Unary { ref arg, .. } => {
                core::slice::from_ref(arg)
            }
            Self::Call { ref args, .. } |
            Self::CallIndirect { ref args, .. } |
            Self::MultiAry { ref args, .. } => {
                args.as_slice(pool)
            }
        }
    }

    /// Get mutable references to the value arguments to this
    /// instruction.
    pub fn arguments_mut<'a>(&'a mut self, pool: &'a mut ir::ValueListPool) -> &mut [Value] {
        match *self {
            Self::DynamicStackLoad { .. } |
            Self::FuncAddr { .. } |
            Self::Jump { .. } |
            Self::NullAry { .. } |
            Self::StackLoad { .. } |
            Self::Trap { .. } |
            Self::UnaryConst { .. } |
            Self::UnaryGlobalValue { .. } |
            Self::UnaryIeee32 { .. } |
            Self::UnaryIeee64 { .. } |
            Self::UnaryImm { .. } => {
                &mut []
            }
            Self::AtomicRmw { args: ref mut args_arity2, .. } |
            Self::Binary { args: ref mut args_arity2, .. } |
            Self::FloatCompare { args: ref mut args_arity2, .. } |
            Self::IntAddTrap { args: ref mut args_arity2, .. } |
            Self::IntCompare { args: ref mut args_arity2, .. } |
            Self::Shuffle { args: ref mut args_arity2, .. } |
            Self::Store { args: ref mut args_arity2, .. } |
            Self::StoreNoOffset { args: ref mut args_arity2, .. } |
            Self::TernaryImm8 { args: ref mut args_arity2, .. } => {
                args_arity2
            }
            Self::AtomicCas { args: ref mut args_arity3, .. } |
            Self::Ternary { args: ref mut args_arity3, .. } => {
                args_arity3
            }
            Self::BinaryImm64 { ref mut arg, .. } |
            Self::BinaryImm8 { ref mut arg, .. } |
            Self::BranchTable { ref mut arg, .. } |
            Self::Brif { ref mut arg, .. } |
            Self::CondTrap { ref mut arg, .. } |
            Self::DynamicStackStore { ref mut arg, .. } |
            Self::IntCompareImm { ref mut arg, .. } |
            Self::Load { ref mut arg, .. } |
            Self::LoadNoOffset { ref mut arg, .. } |
            Self::StackStore { ref mut arg, .. } |
            Self::TableAddr { ref mut arg, .. } |
            Self::Unary { ref mut arg, .. } => {
                core::slice::from_mut(arg)
            }
            Self::Call { ref mut args, .. } |
            Self::CallIndirect { ref mut args, .. } |
            Self::MultiAry { ref mut args, .. } => {
                args.as_mut_slice(pool)
            }
        }
    }

    /// Compare two `InstructionData` for equality.
    ///
    /// This operation requires a reference to a `ValueListPool` to
    /// determine if the contents of any `ValueLists` are equal.
    ///
    /// This operation takes a closure that is allowed to map each
    /// argument value to some other value before the instructions
    /// are compared. This allows various forms of canonicalization.
    pub fn eq<F: Fn(Value) -> Value>(&self, other: &Self, pool: &ir::ValueListPool, mapper: F) -> bool {
        if ::core::mem::discriminant(self) != ::core::mem::discriminant(other) {
            return false;
        }
        match (self, other) {
            (&Self::AtomicCas { opcode: ref opcode1, args: ref args1, flags: ref flags1 }, &Self::AtomicCas { opcode: ref opcode2, args: ref args2, flags: ref flags2 }) => {
                opcode1 == opcode2
                && flags1 == flags2
                && args1.iter().zip(args2.iter()).all(|(a, b)| mapper(*a) == mapper(*b))
            }
            (&Self::AtomicRmw { opcode: ref opcode1, args: ref args1, flags: ref flags1, op: ref op1 }, &Self::AtomicRmw { opcode: ref opcode2, args: ref args2, flags: ref flags2, op: ref op2 }) => {
                opcode1 == opcode2
                && flags1 == flags2
                && op1 == op2
                && args1.iter().zip(args2.iter()).all(|(a, b)| mapper(*a) == mapper(*b))
            }
            (&Self::Binary { opcode: ref opcode1, args: ref args1 }, &Self::Binary { opcode: ref opcode2, args: ref args2 }) => {
                opcode1 == opcode2
                && args1.iter().zip(args2.iter()).all(|(a, b)| mapper(*a) == mapper(*b))
            }
            (&Self::BinaryImm64 { opcode: ref opcode1, arg: ref arg1, imm: ref imm1 }, &Self::BinaryImm64 { opcode: ref opcode2, arg: ref arg2, imm: ref imm2 }) => {
                opcode1 == opcode2
                && imm1 == imm2
                && mapper(*arg1) == mapper(*arg2)
            }
            (&Self::BinaryImm8 { opcode: ref opcode1, arg: ref arg1, imm: ref imm1 }, &Self::BinaryImm8 { opcode: ref opcode2, arg: ref arg2, imm: ref imm2 }) => {
                opcode1 == opcode2
                && imm1 == imm2
                && mapper(*arg1) == mapper(*arg2)
            }
            (&Self::BranchTable { opcode: ref opcode1, arg: ref arg1, table: ref table1 }, &Self::BranchTable { opcode: ref opcode2, arg: ref arg2, table: ref table2 }) => {
                opcode1 == opcode2
                && table1 == table2
                && mapper(*arg1) == mapper(*arg2)
            }
            (&Self::Brif { opcode: ref opcode1, arg: ref arg1, blocks: ref blocks1 }, &Self::Brif { opcode: ref opcode2, arg: ref arg2, blocks: ref blocks2 }) => {
                opcode1 == opcode2
                && mapper(*arg1) == mapper(*arg2)
                && blocks1.iter().zip(blocks2.iter()).all(|(a, b)| a.block(pool) == b.block(pool))
            }
            (&Self::Call { opcode: ref opcode1, args: ref args1, func_ref: ref func_ref1 }, &Self::Call { opcode: ref opcode2, args: ref args2, func_ref: ref func_ref2 }) => {
                opcode1 == opcode2
                && func_ref1 == func_ref2
                && args1.as_slice(pool).iter().zip(args2.as_slice(pool).iter()).all(|(a, b)| mapper(*a) == mapper(*b))
            }
            (&Self::CallIndirect { opcode: ref opcode1, args: ref args1, sig_ref: ref sig_ref1 }, &Self::CallIndirect { opcode: ref opcode2, args: ref args2, sig_ref: ref sig_ref2 }) => {
                opcode1 == opcode2
                && sig_ref1 == sig_ref2
                && args1.as_slice(pool).iter().zip(args2.as_slice(pool).iter()).all(|(a, b)| mapper(*a) == mapper(*b))
            }
            (&Self::CondTrap { opcode: ref opcode1, arg: ref arg1, code: ref code1 }, &Self::CondTrap { opcode: ref opcode2, arg: ref arg2, code: ref code2 }) => {
                opcode1 == opcode2
                && code1 == code2
                && mapper(*arg1) == mapper(*arg2)
            }
            (&Self::DynamicStackLoad { opcode: ref opcode1, dynamic_stack_slot: ref dynamic_stack_slot1 }, &Self::DynamicStackLoad { opcode: ref opcode2, dynamic_stack_slot: ref dynamic_stack_slot2 }) => {
                opcode1 == opcode2
                && dynamic_stack_slot1 == dynamic_stack_slot2
            }
            (&Self::DynamicStackStore { opcode: ref opcode1, arg: ref arg1, dynamic_stack_slot: ref dynamic_stack_slot1 }, &Self::DynamicStackStore { opcode: ref opcode2, arg: ref arg2, dynamic_stack_slot: ref dynamic_stack_slot2 }) => {
                opcode1 == opcode2
                && dynamic_stack_slot1 == dynamic_stack_slot2
                && mapper(*arg1) == mapper(*arg2)
            }
            (&Self::FloatCompare { opcode: ref opcode1, args: ref args1, cond: ref cond1 }, &Self::FloatCompare { opcode: ref opcode2, args: ref args2, cond: ref cond2 }) => {
                opcode1 == opcode2
                && cond1 == cond2
                && args1.iter().zip(args2.iter()).all(|(a, b)| mapper(*a) == mapper(*b))
            }
            (&Self::FuncAddr { opcode: ref opcode1, func_ref: ref func_ref1 }, &Self::FuncAddr { opcode: ref opcode2, func_ref: ref func_ref2 }) => {
                opcode1 == opcode2
                && func_ref1 == func_ref2
            }
            (&Self::IntAddTrap { opcode: ref opcode1, args: ref args1, code: ref code1 }, &Self::IntAddTrap { opcode: ref opcode2, args: ref args2, code: ref code2 }) => {
                opcode1 == opcode2
                && code1 == code2
                && args1.iter().zip(args2.iter()).all(|(a, b)| mapper(*a) == mapper(*b))
            }
            (&Self::IntCompare { opcode: ref opcode1, args: ref args1, cond: ref cond1 }, &Self::IntCompare { opcode: ref opcode2, args: ref args2, cond: ref cond2 }) => {
                opcode1 == opcode2
                && cond1 == cond2
                && args1.iter().zip(args2.iter()).all(|(a, b)| mapper(*a) == mapper(*b))
            }
            (&Self::IntCompareImm { opcode: ref opcode1, arg: ref arg1, cond: ref cond1, imm: ref imm1 }, &Self::IntCompareImm { opcode: ref opcode2, arg: ref arg2, cond: ref cond2, imm: ref imm2 }) => {
                opcode1 == opcode2
                && cond1 == cond2
                && imm1 == imm2
                && mapper(*arg1) == mapper(*arg2)
            }
            (&Self::Jump { opcode: ref opcode1, destination: ref destination1 }, &Self::Jump { opcode: ref opcode2, destination: ref destination2 }) => {
                opcode1 == opcode2
                && destination1 == destination2
            }
            (&Self::Load { opcode: ref opcode1, arg: ref arg1, flags: ref flags1, offset: ref offset1 }, &Self::Load { opcode: ref opcode2, arg: ref arg2, flags: ref flags2, offset: ref offset2 }) => {
                opcode1 == opcode2
                && flags1 == flags2
                && offset1 == offset2
                && mapper(*arg1) == mapper(*arg2)
            }
            (&Self::LoadNoOffset { opcode: ref opcode1, arg: ref arg1, flags: ref flags1 }, &Self::LoadNoOffset { opcode: ref opcode2, arg: ref arg2, flags: ref flags2 }) => {
                opcode1 == opcode2
                && flags1 == flags2
                && mapper(*arg1) == mapper(*arg2)
            }
            (&Self::MultiAry { opcode: ref opcode1, args: ref args1 }, &Self::MultiAry { opcode: ref opcode2, args: ref args2 }) => {
                opcode1 == opcode2
                && args1.as_slice(pool).iter().zip(args2.as_slice(pool).iter()).all(|(a, b)| mapper(*a) == mapper(*b))
            }
            (&Self::NullAry { opcode: ref opcode1 }, &Self::NullAry { opcode: ref opcode2 }) => {
                opcode1 == opcode2
            }
            (&Self::Shuffle { opcode: ref opcode1, args: ref args1, imm: ref imm1 }, &Self::Shuffle { opcode: ref opcode2, args: ref args2, imm: ref imm2 }) => {
                opcode1 == opcode2
                && imm1 == imm2
                && args1.iter().zip(args2.iter()).all(|(a, b)| mapper(*a) == mapper(*b))
            }
            (&Self::StackLoad { opcode: ref opcode1, stack_slot: ref stack_slot1, offset: ref offset1 }, &Self::StackLoad { opcode: ref opcode2, stack_slot: ref stack_slot2, offset: ref offset2 }) => {
                opcode1 == opcode2
                && stack_slot1 == stack_slot2
                && offset1 == offset2
            }
            (&Self::StackStore { opcode: ref opcode1, arg: ref arg1, stack_slot: ref stack_slot1, offset: ref offset1 }, &Self::StackStore { opcode: ref opcode2, arg: ref arg2, stack_slot: ref stack_slot2, offset: ref offset2 }) => {
                opcode1 == opcode2
                && stack_slot1 == stack_slot2
                && offset1 == offset2
                && mapper(*arg1) == mapper(*arg2)
            }
            (&Self::Store { opcode: ref opcode1, args: ref args1, flags: ref flags1, offset: ref offset1 }, &Self::Store { opcode: ref opcode2, args: ref args2, flags: ref flags2, offset: ref offset2 }) => {
                opcode1 == opcode2
                && flags1 == flags2
                && offset1 == offset2
                && args1.iter().zip(args2.iter()).all(|(a, b)| mapper(*a) == mapper(*b))
            }
            (&Self::StoreNoOffset { opcode: ref opcode1, args: ref args1, flags: ref flags1 }, &Self::StoreNoOffset { opcode: ref opcode2, args: ref args2, flags: ref flags2 }) => {
                opcode1 == opcode2
                && flags1 == flags2
                && args1.iter().zip(args2.iter()).all(|(a, b)| mapper(*a) == mapper(*b))
            }
            (&Self::TableAddr { opcode: ref opcode1, arg: ref arg1, table: ref table1, offset: ref offset1 }, &Self::TableAddr { opcode: ref opcode2, arg: ref arg2, table: ref table2, offset: ref offset2 }) => {
                opcode1 == opcode2
                && table1 == table2
                && offset1 == offset2
                && mapper(*arg1) == mapper(*arg2)
            }
            (&Self::Ternary { opcode: ref opcode1, args: ref args1 }, &Self::Ternary { opcode: ref opcode2, args: ref args2 }) => {
                opcode1 == opcode2
                && args1.iter().zip(args2.iter()).all(|(a, b)| mapper(*a) == mapper(*b))
            }
            (&Self::TernaryImm8 { opcode: ref opcode1, args: ref args1, imm: ref imm1 }, &Self::TernaryImm8 { opcode: ref opcode2, args: ref args2, imm: ref imm2 }) => {
                opcode1 == opcode2
                && imm1 == imm2
                && args1.iter().zip(args2.iter()).all(|(a, b)| mapper(*a) == mapper(*b))
            }
            (&Self::Trap { opcode: ref opcode1, code: ref code1 }, &Self::Trap { opcode: ref opcode2, code: ref code2 }) => {
                opcode1 == opcode2
                && code1 == code2
            }
            (&Self::Unary { opcode: ref opcode1, arg: ref arg1 }, &Self::Unary { opcode: ref opcode2, arg: ref arg2 }) => {
                opcode1 == opcode2
                && mapper(*arg1) == mapper(*arg2)
            }
            (&Self::UnaryConst { opcode: ref opcode1, constant_handle: ref constant_handle1 }, &Self::UnaryConst { opcode: ref opcode2, constant_handle: ref constant_handle2 }) => {
                opcode1 == opcode2
                && constant_handle1 == constant_handle2
            }
            (&Self::UnaryGlobalValue { opcode: ref opcode1, global_value: ref global_value1 }, &Self::UnaryGlobalValue { opcode: ref opcode2, global_value: ref global_value2 }) => {
                opcode1 == opcode2
                && global_value1 == global_value2
            }
            (&Self::UnaryIeee32 { opcode: ref opcode1, imm: ref imm1 }, &Self::UnaryIeee32 { opcode: ref opcode2, imm: ref imm2 }) => {
                opcode1 == opcode2
                && imm1 == imm2
            }
            (&Self::UnaryIeee64 { opcode: ref opcode1, imm: ref imm1 }, &Self::UnaryIeee64 { opcode: ref opcode2, imm: ref imm2 }) => {
                opcode1 == opcode2
                && imm1 == imm2
            }
            (&Self::UnaryImm { opcode: ref opcode1, imm: ref imm1 }, &Self::UnaryImm { opcode: ref opcode2, imm: ref imm2 }) => {
                opcode1 == opcode2
                && imm1 == imm2
            }
            _ => unreachable!()
        }
    }

    /// Hash an `InstructionData`.
    ///
    /// This operation requires a reference to a `ValueListPool` to
    /// hash the contents of any `ValueLists`.
    ///
    /// This operation takes a closure that is allowed to map each
    /// argument value to some other value before it is hashed. This
    /// allows various forms of canonicalization.
    pub fn hash<H: ::core::hash::Hasher, F: Fn(Value) -> Value>(&self, state: &mut H, pool: &ir::ValueListPool, mapper: F) {
        match *self {
            Self::AtomicCas{opcode, ref args, flags} => {
                ::core::hash::Hash::hash( &::core::mem::discriminant(self), state);
                ::core::hash::Hash::hash(&opcode, state);
                ::core::hash::Hash::hash(&flags, state);
                ::core::hash::Hash::hash(&args.len(), state);
                for &arg in args {
                    let arg = mapper(arg);
                    ::core::hash::Hash::hash(&arg, state);
                }
            }
            Self::AtomicRmw{opcode, ref args, flags, op} => {
                ::core::hash::Hash::hash( &::core::mem::discriminant(self), state);
                ::core::hash::Hash::hash(&opcode, state);
                ::core::hash::Hash::hash(&flags, state);
                ::core::hash::Hash::hash(&op, state);
                ::core::hash::Hash::hash(&args.len(), state);
                for &arg in args {
                    let arg = mapper(arg);
                    ::core::hash::Hash::hash(&arg, state);
                }
            }
            Self::Binary{opcode, ref args} => {
                ::core::hash::Hash::hash( &::core::mem::discriminant(self), state);
                ::core::hash::Hash::hash(&opcode, state);
                ::core::hash::Hash::hash(&args.len(), state);
                for &arg in args {
                    let arg = mapper(arg);
                    ::core::hash::Hash::hash(&arg, state);
                }
            }
            Self::BinaryImm64{opcode, ref arg, imm} => {
                ::core::hash::Hash::hash( &::core::mem::discriminant(self), state);
                ::core::hash::Hash::hash(&opcode, state);
                ::core::hash::Hash::hash(&imm, state);
                ::core::hash::Hash::hash(&1, state);
                for &arg in std::slice::from_ref(arg) {
                    let arg = mapper(arg);
                    ::core::hash::Hash::hash(&arg, state);
                }
            }
            Self::BinaryImm8{opcode, ref arg, imm} => {
                ::core::hash::Hash::hash( &::core::mem::discriminant(self), state);
                ::core::hash::Hash::hash(&opcode, state);
                ::core::hash::Hash::hash(&imm, state);
                ::core::hash::Hash::hash(&1, state);
                for &arg in std::slice::from_ref(arg) {
                    let arg = mapper(arg);
                    ::core::hash::Hash::hash(&arg, state);
                }
            }
            Self::BranchTable{opcode, ref arg, table} => {
                ::core::hash::Hash::hash( &::core::mem::discriminant(self), state);
                ::core::hash::Hash::hash(&opcode, state);
                ::core::hash::Hash::hash(&table, state);
                ::core::hash::Hash::hash(&1, state);
                for &arg in std::slice::from_ref(arg) {
                    let arg = mapper(arg);
                    ::core::hash::Hash::hash(&arg, state);
                }
            }
            Self::Brif{opcode, ref arg, ref blocks} => {
                ::core::hash::Hash::hash( &::core::mem::discriminant(self), state);
                ::core::hash::Hash::hash(&opcode, state);
                ::core::hash::Hash::hash(&1, state);
                for &arg in std::slice::from_ref(arg) {
                    let arg = mapper(arg);
                    ::core::hash::Hash::hash(&arg, state);
                }
                ::core::hash::Hash::hash(&blocks.len(), state);
                for &block in blocks {
                    ::core::hash::Hash::hash(&block.block(pool), state);
                    for &arg in block.args_slice(pool) {
                        let arg = mapper(arg);
                        ::core::hash::Hash::hash(&arg, state);
                    }
                }
            }
            Self::Call{opcode, ref args, func_ref} => {
                ::core::hash::Hash::hash( &::core::mem::discriminant(self), state);
                ::core::hash::Hash::hash(&opcode, state);
                ::core::hash::Hash::hash(&func_ref, state);
                ::core::hash::Hash::hash(&args.len(pool), state);
                for &arg in args.as_slice(pool) {
                    let arg = mapper(arg);
                    ::core::hash::Hash::hash(&arg, state);
                }
            }
            Self::CallIndirect{opcode, ref args, sig_ref} => {
                ::core::hash::Hash::hash( &::core::mem::discriminant(self), state);
                ::core::hash::Hash::hash(&opcode, state);
                ::core::hash::Hash::hash(&sig_ref, state);
                ::core::hash::Hash::hash(&args.len(pool), state);
                for &arg in args.as_slice(pool) {
                    let arg = mapper(arg);
                    ::core::hash::Hash::hash(&arg, state);
                }
            }
            Self::CondTrap{opcode, ref arg, code} => {
                ::core::hash::Hash::hash( &::core::mem::discriminant(self), state);
                ::core::hash::Hash::hash(&opcode, state);
                ::core::hash::Hash::hash(&code, state);
                ::core::hash::Hash::hash(&1, state);
                for &arg in std::slice::from_ref(arg) {
                    let arg = mapper(arg);
                    ::core::hash::Hash::hash(&arg, state);
                }
            }
            Self::DynamicStackLoad{opcode, dynamic_stack_slot} => {
                ::core::hash::Hash::hash( &::core::mem::discriminant(self), state);
                ::core::hash::Hash::hash(&opcode, state);
                ::core::hash::Hash::hash(&dynamic_stack_slot, state);
                ::core::hash::Hash::hash(&0, state);
                for &arg in &[] {
                    let arg = mapper(arg);
                    ::core::hash::Hash::hash(&arg, state);
                }
            }
            Self::DynamicStackStore{opcode, ref arg, dynamic_stack_slot} => {
                ::core::hash::Hash::hash( &::core::mem::discriminant(self), state);
                ::core::hash::Hash::hash(&opcode, state);
                ::core::hash::Hash::hash(&dynamic_stack_slot, state);
                ::core::hash::Hash::hash(&1, state);
                for &arg in std::slice::from_ref(arg) {
                    let arg = mapper(arg);
                    ::core::hash::Hash::hash(&arg, state);
                }
            }
            Self::FloatCompare{opcode, ref args, cond} => {
                ::core::hash::Hash::hash( &::core::mem::discriminant(self), state);
                ::core::hash::Hash::hash(&opcode, state);
                ::core::hash::Hash::hash(&cond, state);
                ::core::hash::Hash::hash(&args.len(), state);
                for &arg in args {
                    let arg = mapper(arg);
                    ::core::hash::Hash::hash(&arg, state);
                }
            }
            Self::FuncAddr{opcode, func_ref} => {
                ::core::hash::Hash::hash( &::core::mem::discriminant(self), state);
                ::core::hash::Hash::hash(&opcode, state);
                ::core::hash::Hash::hash(&func_ref, state);
                ::core::hash::Hash::hash(&0, state);
                for &arg in &[] {
                    let arg = mapper(arg);
                    ::core::hash::Hash::hash(&arg, state);
                }
            }
            Self::IntAddTrap{opcode, ref args, code} => {
                ::core::hash::Hash::hash( &::core::mem::discriminant(self), state);
                ::core::hash::Hash::hash(&opcode, state);
                ::core::hash::Hash::hash(&code, state);
                ::core::hash::Hash::hash(&args.len(), state);
                for &arg in args {
                    let arg = mapper(arg);
                    ::core::hash::Hash::hash(&arg, state);
                }
            }
            Self::IntCompare{opcode, ref args, cond} => {
                ::core::hash::Hash::hash( &::core::mem::discriminant(self), state);
                ::core::hash::Hash::hash(&opcode, state);
                ::core::hash::Hash::hash(&cond, state);
                ::core::hash::Hash::hash(&args.len(), state);
                for &arg in args {
                    let arg = mapper(arg);
                    ::core::hash::Hash::hash(&arg, state);
                }
            }
            Self::IntCompareImm{opcode, ref arg, cond, imm} => {
                ::core::hash::Hash::hash( &::core::mem::discriminant(self), state);
                ::core::hash::Hash::hash(&opcode, state);
                ::core::hash::Hash::hash(&cond, state);
                ::core::hash::Hash::hash(&imm, state);
                ::core::hash::Hash::hash(&1, state);
                for &arg in std::slice::from_ref(arg) {
                    let arg = mapper(arg);
                    ::core::hash::Hash::hash(&arg, state);
                }
            }
            Self::Jump{opcode, ref destination} => {
                ::core::hash::Hash::hash( &::core::mem::discriminant(self), state);
                ::core::hash::Hash::hash(&opcode, state);
                ::core::hash::Hash::hash(&0, state);
                for &arg in &[] {
                    let arg = mapper(arg);
                    ::core::hash::Hash::hash(&arg, state);
                }
                ::core::hash::Hash::hash(&1, state);
                for &block in std::slice::from_ref(destination) {
                    ::core::hash::Hash::hash(&block.block(pool), state);
                    for &arg in block.args_slice(pool) {
                        let arg = mapper(arg);
                        ::core::hash::Hash::hash(&arg, state);
                    }
                }
            }
            Self::Load{opcode, ref arg, flags, offset} => {
                ::core::hash::Hash::hash( &::core::mem::discriminant(self), state);
                ::core::hash::Hash::hash(&opcode, state);
                ::core::hash::Hash::hash(&flags, state);
                ::core::hash::Hash::hash(&offset, state);
                ::core::hash::Hash::hash(&1, state);
                for &arg in std::slice::from_ref(arg) {
                    let arg = mapper(arg);
                    ::core::hash::Hash::hash(&arg, state);
                }
            }
            Self::LoadNoOffset{opcode, ref arg, flags} => {
                ::core::hash::Hash::hash( &::core::mem::discriminant(self), state);
                ::core::hash::Hash::hash(&opcode, state);
                ::core::hash::Hash::hash(&flags, state);
                ::core::hash::Hash::hash(&1, state);
                for &arg in std::slice::from_ref(arg) {
                    let arg = mapper(arg);
                    ::core::hash::Hash::hash(&arg, state);
                }
            }
            Self::MultiAry{opcode, ref args} => {
                ::core::hash::Hash::hash( &::core::mem::discriminant(self), state);
                ::core::hash::Hash::hash(&opcode, state);
                ::core::hash::Hash::hash(&args.len(pool), state);
                for &arg in args.as_slice(pool) {
                    let arg = mapper(arg);
                    ::core::hash::Hash::hash(&arg, state);
                }
            }
            Self::NullAry{opcode} => {
                ::core::hash::Hash::hash( &::core::mem::discriminant(self), state);
                ::core::hash::Hash::hash(&opcode, state);
                ::core::hash::Hash::hash(&0, state);
                for &arg in &[] {
                    let arg = mapper(arg);
                    ::core::hash::Hash::hash(&arg, state);
                }
            }
            Self::Shuffle{opcode, ref args, imm} => {
                ::core::hash::Hash::hash( &::core::mem::discriminant(self), state);
                ::core::hash::Hash::hash(&opcode, state);
                ::core::hash::Hash::hash(&imm, state);
                ::core::hash::Hash::hash(&args.len(), state);
                for &arg in args {
                    let arg = mapper(arg);
                    ::core::hash::Hash::hash(&arg, state);
                }
            }
            Self::StackLoad{opcode, stack_slot, offset} => {
                ::core::hash::Hash::hash( &::core::mem::discriminant(self), state);
                ::core::hash::Hash::hash(&opcode, state);
                ::core::hash::Hash::hash(&stack_slot, state);
                ::core::hash::Hash::hash(&offset, state);
                ::core::hash::Hash::hash(&0, state);
                for &arg in &[] {
                    let arg = mapper(arg);
                    ::core::hash::Hash::hash(&arg, state);
                }
            }
            Self::StackStore{opcode, ref arg, stack_slot, offset} => {
                ::core::hash::Hash::hash( &::core::mem::discriminant(self), state);
                ::core::hash::Hash::hash(&opcode, state);
                ::core::hash::Hash::hash(&stack_slot, state);
                ::core::hash::Hash::hash(&offset, state);
                ::core::hash::Hash::hash(&1, state);
                for &arg in std::slice::from_ref(arg) {
                    let arg = mapper(arg);
                    ::core::hash::Hash::hash(&arg, state);
                }
            }
            Self::Store{opcode, ref args, flags, offset} => {
                ::core::hash::Hash::hash( &::core::mem::discriminant(self), state);
                ::core::hash::Hash::hash(&opcode, state);
                ::core::hash::Hash::hash(&flags, state);
                ::core::hash::Hash::hash(&offset, state);
                ::core::hash::Hash::hash(&args.len(), state);
                for &arg in args {
                    let arg = mapper(arg);
                    ::core::hash::Hash::hash(&arg, state);
                }
            }
            Self::StoreNoOffset{opcode, ref args, flags} => {
                ::core::hash::Hash::hash( &::core::mem::discriminant(self), state);
                ::core::hash::Hash::hash(&opcode, state);
                ::core::hash::Hash::hash(&flags, state);
                ::core::hash::Hash::hash(&args.len(), state);
                for &arg in args {
                    let arg = mapper(arg);
                    ::core::hash::Hash::hash(&arg, state);
                }
            }
            Self::TableAddr{opcode, ref arg, table, offset} => {
                ::core::hash::Hash::hash( &::core::mem::discriminant(self), state);
                ::core::hash::Hash::hash(&opcode, state);
                ::core::hash::Hash::hash(&table, state);
                ::core::hash::Hash::hash(&offset, state);
                ::core::hash::Hash::hash(&1, state);
                for &arg in std::slice::from_ref(arg) {
                    let arg = mapper(arg);
                    ::core::hash::Hash::hash(&arg, state);
                }
            }
            Self::Ternary{opcode, ref args} => {
                ::core::hash::Hash::hash( &::core::mem::discriminant(self), state);
                ::core::hash::Hash::hash(&opcode, state);
                ::core::hash::Hash::hash(&args.len(), state);
                for &arg in args {
                    let arg = mapper(arg);
                    ::core::hash::Hash::hash(&arg, state);
                }
            }
            Self::TernaryImm8{opcode, ref args, imm} => {
                ::core::hash::Hash::hash( &::core::mem::discriminant(self), state);
                ::core::hash::Hash::hash(&opcode, state);
                ::core::hash::Hash::hash(&imm, state);
                ::core::hash::Hash::hash(&args.len(), state);
                for &arg in args {
                    let arg = mapper(arg);
                    ::core::hash::Hash::hash(&arg, state);
                }
            }
            Self::Trap{opcode, code} => {
                ::core::hash::Hash::hash( &::core::mem::discriminant(self), state);
                ::core::hash::Hash::hash(&opcode, state);
                ::core::hash::Hash::hash(&code, state);
                ::core::hash::Hash::hash(&0, state);
                for &arg in &[] {
                    let arg = mapper(arg);
                    ::core::hash::Hash::hash(&arg, state);
                }
            }
            Self::Unary{opcode, ref arg} => {
                ::core::hash::Hash::hash( &::core::mem::discriminant(self), state);
                ::core::hash::Hash::hash(&opcode, state);
                ::core::hash::Hash::hash(&1, state);
                for &arg in std::slice::from_ref(arg) {
                    let arg = mapper(arg);
                    ::core::hash::Hash::hash(&arg, state);
                }
            }
            Self::UnaryConst{opcode, constant_handle} => {
                ::core::hash::Hash::hash( &::core::mem::discriminant(self), state);
                ::core::hash::Hash::hash(&opcode, state);
                ::core::hash::Hash::hash(&constant_handle, state);
                ::core::hash::Hash::hash(&0, state);
                for &arg in &[] {
                    let arg = mapper(arg);
                    ::core::hash::Hash::hash(&arg, state);
                }
            }
            Self::UnaryGlobalValue{opcode, global_value} => {
                ::core::hash::Hash::hash( &::core::mem::discriminant(self), state);
                ::core::hash::Hash::hash(&opcode, state);
                ::core::hash::Hash::hash(&global_value, state);
                ::core::hash::Hash::hash(&0, state);
                for &arg in &[] {
                    let arg = mapper(arg);
                    ::core::hash::Hash::hash(&arg, state);
                }
            }
            Self::UnaryIeee32{opcode, imm} => {
                ::core::hash::Hash::hash( &::core::mem::discriminant(self), state);
                ::core::hash::Hash::hash(&opcode, state);
                ::core::hash::Hash::hash(&imm, state);
                ::core::hash::Hash::hash(&0, state);
                for &arg in &[] {
                    let arg = mapper(arg);
                    ::core::hash::Hash::hash(&arg, state);
                }
            }
            Self::UnaryIeee64{opcode, imm} => {
                ::core::hash::Hash::hash( &::core::mem::discriminant(self), state);
                ::core::hash::Hash::hash(&opcode, state);
                ::core::hash::Hash::hash(&imm, state);
                ::core::hash::Hash::hash(&0, state);
                for &arg in &[] {
                    let arg = mapper(arg);
                    ::core::hash::Hash::hash(&arg, state);
                }
            }
            Self::UnaryImm{opcode, imm} => {
                ::core::hash::Hash::hash( &::core::mem::discriminant(self), state);
                ::core::hash::Hash::hash(&opcode, state);
                ::core::hash::Hash::hash(&imm, state);
                ::core::hash::Hash::hash(&0, state);
                for &arg in &[] {
                    let arg = mapper(arg);
                    ::core::hash::Hash::hash(&arg, state);
                }
            }
        }
    }

    /// Deep-clone an `InstructionData`, including any referenced lists.
    ///
    /// This operation requires a reference to a `ValueListPool` to
    /// clone the `ValueLists`.
    pub fn deep_clone(&self, pool: &mut ir::ValueListPool) -> Self {
        match *self {
            Self::AtomicCas{opcode, args, flags} => {
                Self::AtomicCas {
                    opcode,
                    args,
                    flags,
                }
            }
            Self::AtomicRmw{opcode, args, flags, op} => {
                Self::AtomicRmw {
                    opcode,
                    args,
                    flags,
                    op,
                }
            }
            Self::Binary{opcode, args} => {
                Self::Binary {
                    opcode,
                    args,
                }
            }
            Self::BinaryImm64{opcode, arg, imm} => {
                Self::BinaryImm64 {
                    opcode,
                    arg,
                    imm,
                }
            }
            Self::BinaryImm8{opcode, arg, imm} => {
                Self::BinaryImm8 {
                    opcode,
                    arg,
                    imm,
                }
            }
            Self::BranchTable{opcode, arg, table} => {
                Self::BranchTable {
                    opcode,
                    arg,
                    table,
                }
            }
            Self::Brif{opcode, arg, blocks} => {
                Self::Brif {
                    opcode,
                    arg,
                    blocks: [blocks[0].deep_clone(pool), blocks[1].deep_clone(pool)],
                }
            }
            Self::Call{opcode, ref args, func_ref} => {
                Self::Call {
                    opcode,
                    args: args.deep_clone(pool),
                    func_ref,
                }
            }
            Self::CallIndirect{opcode, ref args, sig_ref} => {
                Self::CallIndirect {
                    opcode,
                    args: args.deep_clone(pool),
                    sig_ref,
                }
            }
            Self::CondTrap{opcode, arg, code} => {
                Self::CondTrap {
                    opcode,
                    arg,
                    code,
                }
            }
            Self::DynamicStackLoad{opcode, dynamic_stack_slot} => {
                Self::DynamicStackLoad {
                    opcode,
                    dynamic_stack_slot,
                }
            }
            Self::DynamicStackStore{opcode, arg, dynamic_stack_slot} => {
                Self::DynamicStackStore {
                    opcode,
                    arg,
                    dynamic_stack_slot,
                }
            }
            Self::FloatCompare{opcode, args, cond} => {
                Self::FloatCompare {
                    opcode,
                    args,
                    cond,
                }
            }
            Self::FuncAddr{opcode, func_ref} => {
                Self::FuncAddr {
                    opcode,
                    func_ref,
                }
            }
            Self::IntAddTrap{opcode, args, code} => {
                Self::IntAddTrap {
                    opcode,
                    args,
                    code,
                }
            }
            Self::IntCompare{opcode, args, cond} => {
                Self::IntCompare {
                    opcode,
                    args,
                    cond,
                }
            }
            Self::IntCompareImm{opcode, arg, cond, imm} => {
                Self::IntCompareImm {
                    opcode,
                    arg,
                    cond,
                    imm,
                }
            }
            Self::Jump{opcode, destination} => {
                Self::Jump {
                    opcode,
                    destination: destination.deep_clone(pool),
                }
            }
            Self::Load{opcode, arg, flags, offset} => {
                Self::Load {
                    opcode,
                    arg,
                    flags,
                    offset,
                }
            }
            Self::LoadNoOffset{opcode, arg, flags} => {
                Self::LoadNoOffset {
                    opcode,
                    arg,
                    flags,
                }
            }
            Self::MultiAry{opcode, ref args} => {
                Self::MultiAry {
                    opcode,
                    args: args.deep_clone(pool),
                }
            }
            Self::NullAry{opcode} => {
                Self::NullAry {
                    opcode,
                }
            }
            Self::Shuffle{opcode, args, imm} => {
                Self::Shuffle {
                    opcode,
                    args,
                    imm,
                }
            }
            Self::StackLoad{opcode, stack_slot, offset} => {
                Self::StackLoad {
                    opcode,
                    stack_slot,
                    offset,
                }
            }
            Self::StackStore{opcode, arg, stack_slot, offset} => {
                Self::StackStore {
                    opcode,
                    arg,
                    stack_slot,
                    offset,
                }
            }
            Self::Store{opcode, args, flags, offset} => {
                Self::Store {
                    opcode,
                    args,
                    flags,
                    offset,
                }
            }
            Self::StoreNoOffset{opcode, args, flags} => {
                Self::StoreNoOffset {
                    opcode,
                    args,
                    flags,
                }
            }
            Self::TableAddr{opcode, arg, table, offset} => {
                Self::TableAddr {
                    opcode,
                    arg,
                    table,
                    offset,
                }
            }
            Self::Ternary{opcode, args} => {
                Self::Ternary {
                    opcode,
                    args,
                }
            }
            Self::TernaryImm8{opcode, args, imm} => {
                Self::TernaryImm8 {
                    opcode,
                    args,
                    imm,
                }
            }
            Self::Trap{opcode, code} => {
                Self::Trap {
                    opcode,
                    code,
                }
            }
            Self::Unary{opcode, arg} => {
                Self::Unary {
                    opcode,
                    arg,
                }
            }
            Self::UnaryConst{opcode, constant_handle} => {
                Self::UnaryConst {
                    opcode,
                    constant_handle,
                }
            }
            Self::UnaryGlobalValue{opcode, global_value} => {
                Self::UnaryGlobalValue {
                    opcode,
                    global_value,
                }
            }
            Self::UnaryIeee32{opcode, imm} => {
                Self::UnaryIeee32 {
                    opcode,
                    imm,
                }
            }
            Self::UnaryIeee64{opcode, imm} => {
                Self::UnaryIeee64 {
                    opcode,
                    imm,
                }
            }
            Self::UnaryImm{opcode, imm} => {
                Self::UnaryImm {
                    opcode,
                    imm,
                }
            }
        }
    }
}

/// An instruction opcode.
///
/// All instructions from all supported ISAs are present.
#[repr(u8)]
#[derive(Copy, Clone, PartialEq, Eq, Debug, Hash)]
#[cfg_attr(
            feature = "enable-serde",
            derive(serde::Serialize, serde::Deserialize)
        )]
pub enum Opcode {
    /// `jump block_call`. (Jump)
    Jump = 1,
    /// `brif c, block_then, block_else`. (Brif)
    /// Type inferred from `c`.
    Brif,
    /// `br_table x, JT`. (BranchTable)
    BrTable,
    /// `debugtrap`. (NullAry)
    Debugtrap,
    /// `trap code`. (Trap)
    Trap,
    /// `trapz c, code`. (CondTrap)
    /// Type inferred from `c`.
    Trapz,
    /// `resumable_trap code`. (Trap)
    ResumableTrap,
    /// `trapnz c, code`. (CondTrap)
    /// Type inferred from `c`.
    Trapnz,
    /// `resumable_trapnz c, code`. (CondTrap)
    /// Type inferred from `c`.
    ResumableTrapnz,
    /// `return rvals`. (MultiAry)
    Return,
    /// `rvals = call FN, args`. (Call)
    Call,
    /// `rvals = call_indirect SIG, callee, args`. (CallIndirect)
    /// Type inferred from `callee`.
    CallIndirect,
    /// `return_call FN, args`. (Call)
    ReturnCall,
    /// `return_call_indirect SIG, callee, args`. (CallIndirect)
    /// Type inferred from `callee`.
    ReturnCallIndirect,
    /// `addr = func_addr FN`. (FuncAddr)
    FuncAddr,
    /// `a = splat x`. (Unary)
    Splat,
    /// `a = swizzle x, y`. (Binary)
    Swizzle,
    /// `a = insertlane x, y, Idx`. (TernaryImm8)
    /// Type inferred from `x`.
    Insertlane,
    /// `a = extractlane x, Idx`. (BinaryImm8)
    /// Type inferred from `x`.
    Extractlane,
    /// `a = smin x, y`. (Binary)
    /// Type inferred from `x`.
    Smin,
    /// `a = umin x, y`. (Binary)
    /// Type inferred from `x`.
    Umin,
    /// `a = smax x, y`. (Binary)
    /// Type inferred from `x`.
    Smax,
    /// `a = umax x, y`. (Binary)
    /// Type inferred from `x`.
    Umax,
    /// `a = avg_round x, y`. (Binary)
    /// Type inferred from `x`.
    AvgRound,
    /// `a = uadd_sat x, y`. (Binary)
    /// Type inferred from `x`.
    UaddSat,
    /// `a = sadd_sat x, y`. (Binary)
    /// Type inferred from `x`.
    SaddSat,
    /// `a = usub_sat x, y`. (Binary)
    /// Type inferred from `x`.
    UsubSat,
    /// `a = ssub_sat x, y`. (Binary)
    /// Type inferred from `x`.
    SsubSat,
    /// `a = load MemFlags, p, Offset`. (Load)
    Load,
    /// `store MemFlags, x, p, Offset`. (Store)
    /// Type inferred from `x`.
    Store,
    /// `a = uload8 MemFlags, p, Offset`. (Load)
    Uload8,
    /// `a = sload8 MemFlags, p, Offset`. (Load)
    Sload8,
    /// `istore8 MemFlags, x, p, Offset`. (Store)
    /// Type inferred from `x`.
    Istore8,
    /// `a = uload16 MemFlags, p, Offset`. (Load)
    Uload16,
    /// `a = sload16 MemFlags, p, Offset`. (Load)
    Sload16,
    /// `istore16 MemFlags, x, p, Offset`. (Store)
    /// Type inferred from `x`.
    Istore16,
    /// `a = uload32 MemFlags, p, Offset`. (Load)
    /// Type inferred from `p`.
    Uload32,
    /// `a = sload32 MemFlags, p, Offset`. (Load)
    /// Type inferred from `p`.
    Sload32,
    /// `istore32 MemFlags, x, p, Offset`. (Store)
    /// Type inferred from `x`.
    Istore32,
    /// `a = uload8x8 MemFlags, p, Offset`. (Load)
    /// Type inferred from `p`.
    Uload8x8,
    /// `a = sload8x8 MemFlags, p, Offset`. (Load)
    /// Type inferred from `p`.
    Sload8x8,
    /// `a = uload16x4 MemFlags, p, Offset`. (Load)
    /// Type inferred from `p`.
    Uload16x4,
    /// `a = sload16x4 MemFlags, p, Offset`. (Load)
    /// Type inferred from `p`.
    Sload16x4,
    /// `a = uload32x2 MemFlags, p, Offset`. (Load)
    /// Type inferred from `p`.
    Uload32x2,
    /// `a = sload32x2 MemFlags, p, Offset`. (Load)
    /// Type inferred from `p`.
    Sload32x2,
    /// `a = stack_load SS, Offset`. (StackLoad)
    StackLoad,
    /// `stack_store x, SS, Offset`. (StackStore)
    /// Type inferred from `x`.
    StackStore,
    /// `addr = stack_addr SS, Offset`. (StackLoad)
    StackAddr,
    /// `a = dynamic_stack_load DSS`. (DynamicStackLoad)
    DynamicStackLoad,
    /// `dynamic_stack_store x, DSS`. (DynamicStackStore)
    /// Type inferred from `x`.
    DynamicStackStore,
    /// `addr = dynamic_stack_addr DSS`. (DynamicStackLoad)
    DynamicStackAddr,
    /// `a = global_value GV`. (UnaryGlobalValue)
    GlobalValue,
    /// `a = symbol_value GV`. (UnaryGlobalValue)
    SymbolValue,
    /// `a = tls_value GV`. (UnaryGlobalValue)
    TlsValue,
    /// `addr = get_pinned_reg`. (NullAry)
    GetPinnedReg,
    /// `set_pinned_reg addr`. (Unary)
    /// Type inferred from `addr`.
    SetPinnedReg,
    /// `addr = get_frame_pointer`. (NullAry)
    GetFramePointer,
    /// `addr = get_stack_pointer`. (NullAry)
    GetStackPointer,
    /// `addr = get_return_address`. (NullAry)
    GetReturnAddress,
    /// `addr = table_addr T, p, Offset`. (TableAddr)
    TableAddr,
    /// `a = iconst N`. (UnaryImm)
    Iconst,
    /// `a = f32const N`. (UnaryIeee32)
    F32const,
    /// `a = f64const N`. (UnaryIeee64)
    F64const,
    /// `a = vconst N`. (UnaryConst)
    Vconst,
    /// `a = shuffle a, b, mask`. (Shuffle)
    Shuffle,
    /// `a = null`. (NullAry)
    Null,
    /// `nop`. (NullAry)
    Nop,
    /// `a = select c, x, y`. (Ternary)
    /// Type inferred from `x`.
    Select,
    /// `a = select_spectre_guard c, x, y`. (Ternary)
    /// Type inferred from `x`.
    SelectSpectreGuard,
    /// `a = bitselect c, x, y`. (Ternary)
    /// Type inferred from `x`.
    Bitselect,
    /// `a = vselect c, x, y`. (Ternary)
    /// Type inferred from `x`.
    Vselect,
    /// `s = vany_true a`. (Unary)
    /// Type inferred from `a`.
    VanyTrue,
    /// `s = vall_true a`. (Unary)
    /// Type inferred from `a`.
    VallTrue,
    /// `x = vhigh_bits a`. (Unary)
    VhighBits,
    /// `a = icmp Cond, x, y`. (IntCompare)
    /// Type inferred from `x`.
    Icmp,
    /// `a = icmp_imm Cond, x, Y`. (IntCompareImm)
    /// Type inferred from `x`.
    IcmpImm,
    /// `a = iadd x, y`. (Binary)
    /// Type inferred from `x`.
    Iadd,
    /// `a = isub x, y`. (Binary)
    /// Type inferred from `x`.
    Isub,
    /// `a = ineg x`. (Unary)
    /// Type inferred from `x`.
    Ineg,
    /// `a = iabs x`. (Unary)
    /// Type inferred from `x`.
    Iabs,
    /// `a = imul x, y`. (Binary)
    /// Type inferred from `x`.
    Imul,
    /// `a = umulhi x, y`. (Binary)
    /// Type inferred from `x`.
    Umulhi,
    /// `a = smulhi x, y`. (Binary)
    /// Type inferred from `x`.
    Smulhi,
    /// `a = sqmul_round_sat x, y`. (Binary)
    /// Type inferred from `x`.
    SqmulRoundSat,
    /// `a = udiv x, y`. (Binary)
    /// Type inferred from `x`.
    Udiv,
    /// `a = sdiv x, y`. (Binary)
    /// Type inferred from `x`.
    Sdiv,
    /// `a = urem x, y`. (Binary)
    /// Type inferred from `x`.
    Urem,
    /// `a = srem x, y`. (Binary)
    /// Type inferred from `x`.
    Srem,
    /// `a = iadd_imm x, Y`. (BinaryImm64)
    /// Type inferred from `x`.
    IaddImm,
    /// `a = imul_imm x, Y`. (BinaryImm64)
    /// Type inferred from `x`.
    ImulImm,
    /// `a = udiv_imm x, Y`. (BinaryImm64)
    /// Type inferred from `x`.
    UdivImm,
    /// `a = sdiv_imm x, Y`. (BinaryImm64)
    /// Type inferred from `x`.
    SdivImm,
    /// `a = urem_imm x, Y`. (BinaryImm64)
    /// Type inferred from `x`.
    UremImm,
    /// `a = srem_imm x, Y`. (BinaryImm64)
    /// Type inferred from `x`.
    SremImm,
    /// `a = irsub_imm x, Y`. (BinaryImm64)
    /// Type inferred from `x`.
    IrsubImm,
    /// `a = iadd_cin x, y, c_in`. (Ternary)
    /// Type inferred from `y`.
    IaddCin,
    /// `a, c_out = iadd_cout x, y`. (Binary)
    /// Type inferred from `x`.
    IaddCout,
    /// `a, c_out = iadd_carry x, y, c_in`. (Ternary)
    /// Type inferred from `y`.
    IaddCarry,
    /// `a = uadd_overflow_trap x, y, code`. (IntAddTrap)
    /// Type inferred from `x`.
    UaddOverflowTrap,
    /// `a = isub_bin x, y, b_in`. (Ternary)
    /// Type inferred from `y`.
    IsubBin,
    /// `a, b_out = isub_bout x, y`. (Binary)
    /// Type inferred from `x`.
    IsubBout,
    /// `a, b_out = isub_borrow x, y, b_in`. (Ternary)
    /// Type inferred from `y`.
    IsubBorrow,
    /// `a = band x, y`. (Binary)
    /// Type inferred from `x`.
    Band,
    /// `a = bor x, y`. (Binary)
    /// Type inferred from `x`.
    Bor,
    /// `a = bxor x, y`. (Binary)
    /// Type inferred from `x`.
    Bxor,
    /// `a = bnot x`. (Unary)
    /// Type inferred from `x`.
    Bnot,
    /// `a = band_not x, y`. (Binary)
    /// Type inferred from `x`.
    BandNot,
    /// `a = bor_not x, y`. (Binary)
    /// Type inferred from `x`.
    BorNot,
    /// `a = bxor_not x, y`. (Binary)
    /// Type inferred from `x`.
    BxorNot,
    /// `a = band_imm x, Y`. (BinaryImm64)
    /// Type inferred from `x`.
    BandImm,
    /// `a = bor_imm x, Y`. (BinaryImm64)
    /// Type inferred from `x`.
    BorImm,
    /// `a = bxor_imm x, Y`. (BinaryImm64)
    /// Type inferred from `x`.
    BxorImm,
    /// `a = rotl x, y`. (Binary)
    /// Type inferred from `x`.
    Rotl,
    /// `a = rotr x, y`. (Binary)
    /// Type inferred from `x`.
    Rotr,
    /// `a = rotl_imm x, Y`. (BinaryImm64)
    /// Type inferred from `x`.
    RotlImm,
    /// `a = rotr_imm x, Y`. (BinaryImm64)
    /// Type inferred from `x`.
    RotrImm,
    /// `a = ishl x, y`. (Binary)
    /// Type inferred from `x`.
    Ishl,
    /// `a = ushr x, y`. (Binary)
    /// Type inferred from `x`.
    Ushr,
    /// `a = sshr x, y`. (Binary)
    /// Type inferred from `x`.
    Sshr,
    /// `a = ishl_imm x, Y`. (BinaryImm64)
    /// Type inferred from `x`.
    IshlImm,
    /// `a = ushr_imm x, Y`. (BinaryImm64)
    /// Type inferred from `x`.
    UshrImm,
    /// `a = sshr_imm x, Y`. (BinaryImm64)
    /// Type inferred from `x`.
    SshrImm,
    /// `a = bitrev x`. (Unary)
    /// Type inferred from `x`.
    Bitrev,
    /// `a = clz x`. (Unary)
    /// Type inferred from `x`.
    Clz,
    /// `a = cls x`. (Unary)
    /// Type inferred from `x`.
    Cls,
    /// `a = ctz x`. (Unary)
    /// Type inferred from `x`.
    Ctz,
    /// `a = bswap x`. (Unary)
    /// Type inferred from `x`.
    Bswap,
    /// `a = popcnt x`. (Unary)
    /// Type inferred from `x`.
    Popcnt,
    /// `a = fcmp Cond, x, y`. (FloatCompare)
    /// Type inferred from `x`.
    Fcmp,
    /// `a = fadd x, y`. (Binary)
    /// Type inferred from `x`.
    Fadd,
    /// `a = fsub x, y`. (Binary)
    /// Type inferred from `x`.
    Fsub,
    /// `a = fmul x, y`. (Binary)
    /// Type inferred from `x`.
    Fmul,
    /// `a = fdiv x, y`. (Binary)
    /// Type inferred from `x`.
    Fdiv,
    /// `a = sqrt x`. (Unary)
    /// Type inferred from `x`.
    Sqrt,
    /// `a = fma x, y, z`. (Ternary)
    /// Type inferred from `y`.
    Fma,
    /// `a = fneg x`. (Unary)
    /// Type inferred from `x`.
    Fneg,
    /// `a = fabs x`. (Unary)
    /// Type inferred from `x`.
    Fabs,
    /// `a = fcopysign x, y`. (Binary)
    /// Type inferred from `x`.
    Fcopysign,
    /// `a = fmin x, y`. (Binary)
    /// Type inferred from `x`.
    Fmin,
    /// `a = fmin_pseudo x, y`. (Binary)
    /// Type inferred from `x`.
    FminPseudo,
    /// `a = fmax x, y`. (Binary)
    /// Type inferred from `x`.
    Fmax,
    /// `a = fmax_pseudo x, y`. (Binary)
    /// Type inferred from `x`.
    FmaxPseudo,
    /// `a = ceil x`. (Unary)
    /// Type inferred from `x`.
    Ceil,
    /// `a = floor x`. (Unary)
    /// Type inferred from `x`.
    Floor,
    /// `a = trunc x`. (Unary)
    /// Type inferred from `x`.
    Trunc,
    /// `a = nearest x`. (Unary)
    /// Type inferred from `x`.
    Nearest,
    /// `a = is_null x`. (Unary)
    /// Type inferred from `x`.
    IsNull,
    /// `a = is_invalid x`. (Unary)
    /// Type inferred from `x`.
    IsInvalid,
    /// `a = bitcast MemFlags, x`. (LoadNoOffset)
    Bitcast,
    /// `a = scalar_to_vector s`. (Unary)
    ScalarToVector,
    /// `a = bmask x`. (Unary)
    Bmask,
    /// `a = ireduce x`. (Unary)
    Ireduce,
    /// `a = snarrow x, y`. (Binary)
    /// Type inferred from `x`.
    Snarrow,
    /// `a = unarrow x, y`. (Binary)
    /// Type inferred from `x`.
    Unarrow,
    /// `a = uunarrow x, y`. (Binary)
    /// Type inferred from `x`.
    Uunarrow,
    /// `a = swiden_low x`. (Unary)
    /// Type inferred from `x`.
    SwidenLow,
    /// `a = swiden_high x`. (Unary)
    /// Type inferred from `x`.
    SwidenHigh,
    /// `a = uwiden_low x`. (Unary)
    /// Type inferred from `x`.
    UwidenLow,
    /// `a = uwiden_high x`. (Unary)
    /// Type inferred from `x`.
    UwidenHigh,
    /// `a = iadd_pairwise x, y`. (Binary)
    /// Type inferred from `x`.
    IaddPairwise,
    /// `a = uextend x`. (Unary)
    Uextend,
    /// `a = sextend x`. (Unary)
    Sextend,
    /// `a = fpromote x`. (Unary)
    Fpromote,
    /// `a = fdemote x`. (Unary)
    Fdemote,
    /// `a = fvdemote x`. (Unary)
    Fvdemote,
    /// `x = fvpromote_low a`. (Unary)
    FvpromoteLow,
    /// `a = fcvt_to_uint x`. (Unary)
    FcvtToUint,
    /// `a = fcvt_to_sint x`. (Unary)
    FcvtToSint,
    /// `a = fcvt_to_uint_sat x`. (Unary)
    FcvtToUintSat,
    /// `a = fcvt_to_sint_sat x`. (Unary)
    FcvtToSintSat,
    /// `a = fcvt_from_uint x`. (Unary)
    FcvtFromUint,
    /// `a = fcvt_from_sint x`. (Unary)
    FcvtFromSint,
    /// `a = fcvt_low_from_sint x`. (Unary)
    FcvtLowFromSint,
    /// `lo, hi = isplit x`. (Unary)
    /// Type inferred from `x`.
    Isplit,
    /// `a = iconcat lo, hi`. (Binary)
    /// Type inferred from `lo`.
    Iconcat,
    /// `a = atomic_rmw MemFlags, AtomicRmwOp, p, x`. (AtomicRmw)
    AtomicRmw,
    /// `a = atomic_cas MemFlags, p, e, x`. (AtomicCas)
    /// Type inferred from `x`.
    AtomicCas,
    /// `a = atomic_load MemFlags, p`. (LoadNoOffset)
    AtomicLoad,
    /// `atomic_store MemFlags, x, p`. (StoreNoOffset)
    /// Type inferred from `x`.
    AtomicStore,
    /// `fence`. (NullAry)
    Fence,
    /// `a = extract_vector x, y`. (BinaryImm8)
    /// Type inferred from `x`.
    ExtractVector,
}

impl Opcode {
    /// True for instructions that terminate the block
    pub fn is_terminator(self) -> bool {
        match self {
            Self::BrTable |
            Self::Brif |
            Self::Jump |
            Self::Return |
            Self::ReturnCall |
            Self::ReturnCallIndirect |
            Self::Trap => {
                true
            }
            _ => {
                false
            }
        }
    }

    /// True for all branch or jump instructions.
    pub fn is_branch(self) -> bool {
        match self {
            Self::BrTable |
            Self::Brif |
            Self::Jump => {
                true
            }
            _ => {
                false
            }
        }
    }

    /// Is this a call instruction?
    pub fn is_call(self) -> bool {
        match self {
            Self::Call |
            Self::CallIndirect |
            Self::ReturnCall |
            Self::ReturnCallIndirect => {
                true
            }
            _ => {
                false
            }
        }
    }

    /// Is this a return instruction?
    pub fn is_return(self) -> bool {
        match self {
            Self::Return |
            Self::ReturnCall |
            Self::ReturnCallIndirect => {
                true
            }
            _ => {
                false
            }
        }
    }

    /// Can this instruction read from memory?
    pub fn can_load(self) -> bool {
        match self {
            Self::AtomicCas |
            Self::AtomicLoad |
            Self::AtomicRmw |
            Self::Debugtrap |
            Self::DynamicStackLoad |
            Self::Load |
            Self::Sload16 |
            Self::Sload16x4 |
            Self::Sload32 |
            Self::Sload32x2 |
            Self::Sload8 |
            Self::Sload8x8 |
            Self::StackLoad |
            Self::Uload16 |
            Self::Uload16x4 |
            Self::Uload32 |
            Self::Uload32x2 |
            Self::Uload8 |
            Self::Uload8x8 => {
                true
            }
            _ => {
                false
            }
        }
    }

    /// Can this instruction write to memory?
    pub fn can_store(self) -> bool {
        match self {
            Self::AtomicCas |
            Self::AtomicRmw |
            Self::AtomicStore |
            Self::Debugtrap |
            Self::DynamicStackStore |
            Self::Istore16 |
            Self::Istore32 |
            Self::Istore8 |
            Self::StackStore |
            Self::Store => {
                true
            }
            _ => {
                false
            }
        }
    }

    /// Can this instruction cause a trap?
    pub fn can_trap(self) -> bool {
        match self {
            Self::FcvtToSint |
            Self::FcvtToUint |
            Self::ResumableTrap |
            Self::ResumableTrapnz |
            Self::Sdiv |
            Self::Srem |
            Self::Trap |
            Self::Trapnz |
            Self::Trapz |
            Self::UaddOverflowTrap |
            Self::Udiv |
            Self::Urem => {
                true
            }
            _ => {
                false
            }
        }
    }

    /// Does this instruction have other side effects besides can_* flags?
    pub fn other_side_effects(self) -> bool {
        match self {
            Self::AtomicCas |
            Self::AtomicLoad |
            Self::AtomicRmw |
            Self::AtomicStore |
            Self::Debugtrap |
            Self::Fence |
            Self::GetPinnedReg |
            Self::SelectSpectreGuard |
            Self::SetPinnedReg => {
                true
            }
            _ => {
                false
            }
        }
    }

    /// Despite having side effects, is this instruction okay to GVN?
    pub fn side_effects_idempotent(self) -> bool {
        match self {
            Self::FcvtToSint |
            Self::FcvtToUint |
            Self::Sdiv |
            Self::SelectSpectreGuard |
            Self::Srem |
            Self::UaddOverflowTrap |
            Self::Udiv |
            Self::Urem => {
                true
            }
            _ => {
                false
            }
        }
    }

}

const OPCODE_FORMAT: [InstructionFormat; 181] = [
    InstructionFormat::Jump, // jump
    InstructionFormat::Brif, // brif
    InstructionFormat::BranchTable, // br_table
    InstructionFormat::NullAry, // debugtrap
    InstructionFormat::Trap, // trap
    InstructionFormat::CondTrap, // trapz
    InstructionFormat::Trap, // resumable_trap
    InstructionFormat::CondTrap, // trapnz
    InstructionFormat::CondTrap, // resumable_trapnz
    InstructionFormat::MultiAry, // return
    InstructionFormat::Call, // call
    InstructionFormat::CallIndirect, // call_indirect
    InstructionFormat::Call, // return_call
    InstructionFormat::CallIndirect, // return_call_indirect
    InstructionFormat::FuncAddr, // func_addr
    InstructionFormat::Unary, // splat
    InstructionFormat::Binary, // swizzle
    InstructionFormat::TernaryImm8, // insertlane
    InstructionFormat::BinaryImm8, // extractlane
    InstructionFormat::Binary, // smin
    InstructionFormat::Binary, // umin
    InstructionFormat::Binary, // smax
    InstructionFormat::Binary, // umax
    InstructionFormat::Binary, // avg_round
    InstructionFormat::Binary, // uadd_sat
    InstructionFormat::Binary, // sadd_sat
    InstructionFormat::Binary, // usub_sat
    InstructionFormat::Binary, // ssub_sat
    InstructionFormat::Load, // load
    InstructionFormat::Store, // store
    InstructionFormat::Load, // uload8
    InstructionFormat::Load, // sload8
    InstructionFormat::Store, // istore8
    InstructionFormat::Load, // uload16
    InstructionFormat::Load, // sload16
    InstructionFormat::Store, // istore16
    InstructionFormat::Load, // uload32
    InstructionFormat::Load, // sload32
    InstructionFormat::Store, // istore32
    InstructionFormat::Load, // uload8x8
    InstructionFormat::Load, // sload8x8
    InstructionFormat::Load, // uload16x4
    InstructionFormat::Load, // sload16x4
    InstructionFormat::Load, // uload32x2
    InstructionFormat::Load, // sload32x2
    InstructionFormat::StackLoad, // stack_load
    InstructionFormat::StackStore, // stack_store
    InstructionFormat::StackLoad, // stack_addr
    InstructionFormat::DynamicStackLoad, // dynamic_stack_load
    InstructionFormat::DynamicStackStore, // dynamic_stack_store
    InstructionFormat::DynamicStackLoad, // dynamic_stack_addr
    InstructionFormat::UnaryGlobalValue, // global_value
    InstructionFormat::UnaryGlobalValue, // symbol_value
    InstructionFormat::UnaryGlobalValue, // tls_value
    InstructionFormat::NullAry, // get_pinned_reg
    InstructionFormat::Unary, // set_pinned_reg
    InstructionFormat::NullAry, // get_frame_pointer
    InstructionFormat::NullAry, // get_stack_pointer
    InstructionFormat::NullAry, // get_return_address
    InstructionFormat::TableAddr, // table_addr
    InstructionFormat::UnaryImm, // iconst
    InstructionFormat::UnaryIeee32, // f32const
    InstructionFormat::UnaryIeee64, // f64const
    InstructionFormat::UnaryConst, // vconst
    InstructionFormat::Shuffle, // shuffle
    InstructionFormat::NullAry, // null
    InstructionFormat::NullAry, // nop
    InstructionFormat::Ternary, // select
    InstructionFormat::Ternary, // select_spectre_guard
    InstructionFormat::Ternary, // bitselect
    InstructionFormat::Ternary, // vselect
    InstructionFormat::Unary, // vany_true
    InstructionFormat::Unary, // vall_true
    InstructionFormat::Unary, // vhigh_bits
    InstructionFormat::IntCompare, // icmp
    InstructionFormat::IntCompareImm, // icmp_imm
    InstructionFormat::Binary, // iadd
    InstructionFormat::Binary, // isub
    InstructionFormat::Unary, // ineg
    InstructionFormat::Unary, // iabs
    InstructionFormat::Binary, // imul
    InstructionFormat::Binary, // umulhi
    InstructionFormat::Binary, // smulhi
    InstructionFormat::Binary, // sqmul_round_sat
    InstructionFormat::Binary, // udiv
    InstructionFormat::Binary, // sdiv
    InstructionFormat::Binary, // urem
    InstructionFormat::Binary, // srem
    InstructionFormat::BinaryImm64, // iadd_imm
    InstructionFormat::BinaryImm64, // imul_imm
    InstructionFormat::BinaryImm64, // udiv_imm
    InstructionFormat::BinaryImm64, // sdiv_imm
    InstructionFormat::BinaryImm64, // urem_imm
    InstructionFormat::BinaryImm64, // srem_imm
    InstructionFormat::BinaryImm64, // irsub_imm
    InstructionFormat::Ternary, // iadd_cin
    InstructionFormat::Binary, // iadd_cout
    InstructionFormat::Ternary, // iadd_carry
    InstructionFormat::IntAddTrap, // uadd_overflow_trap
    InstructionFormat::Ternary, // isub_bin
    InstructionFormat::Binary, // isub_bout
    InstructionFormat::Ternary, // isub_borrow
    InstructionFormat::Binary, // band
    InstructionFormat::Binary, // bor
    InstructionFormat::Binary, // bxor
    InstructionFormat::Unary, // bnot
    InstructionFormat::Binary, // band_not
    InstructionFormat::Binary, // bor_not
    InstructionFormat::Binary, // bxor_not
    InstructionFormat::BinaryImm64, // band_imm
    InstructionFormat::BinaryImm64, // bor_imm
    InstructionFormat::BinaryImm64, // bxor_imm
    InstructionFormat::Binary, // rotl
    InstructionFormat::Binary, // rotr
    InstructionFormat::BinaryImm64, // rotl_imm
    InstructionFormat::BinaryImm64, // rotr_imm
    InstructionFormat::Binary, // ishl
    InstructionFormat::Binary, // ushr
    InstructionFormat::Binary, // sshr
    InstructionFormat::BinaryImm64, // ishl_imm
    InstructionFormat::BinaryImm64, // ushr_imm
    InstructionFormat::BinaryImm64, // sshr_imm
    InstructionFormat::Unary, // bitrev
    InstructionFormat::Unary, // clz
    InstructionFormat::Unary, // cls
    InstructionFormat::Unary, // ctz
    InstructionFormat::Unary, // bswap
    InstructionFormat::Unary, // popcnt
    InstructionFormat::FloatCompare, // fcmp
    InstructionFormat::Binary, // fadd
    InstructionFormat::Binary, // fsub
    InstructionFormat::Binary, // fmul
    InstructionFormat::Binary, // fdiv
    InstructionFormat::Unary, // sqrt
    InstructionFormat::Ternary, // fma
    InstructionFormat::Unary, // fneg
    InstructionFormat::Unary, // fabs
    InstructionFormat::Binary, // fcopysign
    InstructionFormat::Binary, // fmin
    InstructionFormat::Binary, // fmin_pseudo
    InstructionFormat::Binary, // fmax
    InstructionFormat::Binary, // fmax_pseudo
    InstructionFormat::Unary, // ceil
    InstructionFormat::Unary, // floor
    InstructionFormat::Unary, // trunc
    InstructionFormat::Unary, // nearest
    InstructionFormat::Unary, // is_null
    InstructionFormat::Unary, // is_invalid
    InstructionFormat::LoadNoOffset, // bitcast
    InstructionFormat::Unary, // scalar_to_vector
    InstructionFormat::Unary, // bmask
    InstructionFormat::Unary, // ireduce
    InstructionFormat::Binary, // snarrow
    InstructionFormat::Binary, // unarrow
    InstructionFormat::Binary, // uunarrow
    InstructionFormat::Unary, // swiden_low
    InstructionFormat::Unary, // swiden_high
    InstructionFormat::Unary, // uwiden_low
    InstructionFormat::Unary, // uwiden_high
    InstructionFormat::Binary, // iadd_pairwise
    InstructionFormat::Unary, // uextend
    InstructionFormat::Unary, // sextend
    InstructionFormat::Unary, // fpromote
    InstructionFormat::Unary, // fdemote
    InstructionFormat::Unary, // fvdemote
    InstructionFormat::Unary, // fvpromote_low
    InstructionFormat::Unary, // fcvt_to_uint
    InstructionFormat::Unary, // fcvt_to_sint
    InstructionFormat::Unary, // fcvt_to_uint_sat
    InstructionFormat::Unary, // fcvt_to_sint_sat
    InstructionFormat::Unary, // fcvt_from_uint
    InstructionFormat::Unary, // fcvt_from_sint
    InstructionFormat::Unary, // fcvt_low_from_sint
    InstructionFormat::Unary, // isplit
    InstructionFormat::Binary, // iconcat
    InstructionFormat::AtomicRmw, // atomic_rmw
    InstructionFormat::AtomicCas, // atomic_cas
    InstructionFormat::LoadNoOffset, // atomic_load
    InstructionFormat::StoreNoOffset, // atomic_store
    InstructionFormat::NullAry, // fence
    InstructionFormat::BinaryImm8, // extract_vector
];

fn opcode_name(opc: Opcode) -> &'static str {
    match opc {
        Opcode::AtomicCas => {
            "atomic_cas"
        }
        Opcode::AtomicLoad => {
            "atomic_load"
        }
        Opcode::AtomicRmw => {
            "atomic_rmw"
        }
        Opcode::AtomicStore => {
            "atomic_store"
        }
        Opcode::AvgRound => {
            "avg_round"
        }
        Opcode::Band => {
            "band"
        }
        Opcode::BandImm => {
            "band_imm"
        }
        Opcode::BandNot => {
            "band_not"
        }
        Opcode::Bitcast => {
            "bitcast"
        }
        Opcode::Bitrev => {
            "bitrev"
        }
        Opcode::Bitselect => {
            "bitselect"
        }
        Opcode::Bmask => {
            "bmask"
        }
        Opcode::Bnot => {
            "bnot"
        }
        Opcode::Bor => {
            "bor"
        }
        Opcode::BorImm => {
            "bor_imm"
        }
        Opcode::BorNot => {
            "bor_not"
        }
        Opcode::BrTable => {
            "br_table"
        }
        Opcode::Brif => {
            "brif"
        }
        Opcode::Bswap => {
            "bswap"
        }
        Opcode::Bxor => {
            "bxor"
        }
        Opcode::BxorImm => {
            "bxor_imm"
        }
        Opcode::BxorNot => {
            "bxor_not"
        }
        Opcode::Call => {
            "call"
        }
        Opcode::CallIndirect => {
            "call_indirect"
        }
        Opcode::Ceil => {
            "ceil"
        }
        Opcode::Cls => {
            "cls"
        }
        Opcode::Clz => {
            "clz"
        }
        Opcode::Ctz => {
            "ctz"
        }
        Opcode::Debugtrap => {
            "debugtrap"
        }
        Opcode::DynamicStackAddr => {
            "dynamic_stack_addr"
        }
        Opcode::DynamicStackLoad => {
            "dynamic_stack_load"
        }
        Opcode::DynamicStackStore => {
            "dynamic_stack_store"
        }
        Opcode::ExtractVector => {
            "extract_vector"
        }
        Opcode::Extractlane => {
            "extractlane"
        }
        Opcode::F32const => {
            "f32const"
        }
        Opcode::F64const => {
            "f64const"
        }
        Opcode::Fabs => {
            "fabs"
        }
        Opcode::Fadd => {
            "fadd"
        }
        Opcode::Fcmp => {
            "fcmp"
        }
        Opcode::Fcopysign => {
            "fcopysign"
        }
        Opcode::FcvtFromSint => {
            "fcvt_from_sint"
        }
        Opcode::FcvtFromUint => {
            "fcvt_from_uint"
        }
        Opcode::FcvtLowFromSint => {
            "fcvt_low_from_sint"
        }
        Opcode::FcvtToSint => {
            "fcvt_to_sint"
        }
        Opcode::FcvtToSintSat => {
            "fcvt_to_sint_sat"
        }
        Opcode::FcvtToUint => {
            "fcvt_to_uint"
        }
        Opcode::FcvtToUintSat => {
            "fcvt_to_uint_sat"
        }
        Opcode::Fdemote => {
            "fdemote"
        }
        Opcode::Fdiv => {
            "fdiv"
        }
        Opcode::Fence => {
            "fence"
        }
        Opcode::Floor => {
            "floor"
        }
        Opcode::Fma => {
            "fma"
        }
        Opcode::Fmax => {
            "fmax"
        }
        Opcode::FmaxPseudo => {
            "fmax_pseudo"
        }
        Opcode::Fmin => {
            "fmin"
        }
        Opcode::FminPseudo => {
            "fmin_pseudo"
        }
        Opcode::Fmul => {
            "fmul"
        }
        Opcode::Fneg => {
            "fneg"
        }
        Opcode::Fpromote => {
            "fpromote"
        }
        Opcode::Fsub => {
            "fsub"
        }
        Opcode::FuncAddr => {
            "func_addr"
        }
        Opcode::Fvdemote => {
            "fvdemote"
        }
        Opcode::FvpromoteLow => {
            "fvpromote_low"
        }
        Opcode::GetFramePointer => {
            "get_frame_pointer"
        }
        Opcode::GetPinnedReg => {
            "get_pinned_reg"
        }
        Opcode::GetReturnAddress => {
            "get_return_address"
        }
        Opcode::GetStackPointer => {
            "get_stack_pointer"
        }
        Opcode::GlobalValue => {
            "global_value"
        }
        Opcode::Iabs => {
            "iabs"
        }
        Opcode::Iadd => {
            "iadd"
        }
        Opcode::IaddCarry => {
            "iadd_carry"
        }
        Opcode::IaddCin => {
            "iadd_cin"
        }
        Opcode::IaddCout => {
            "iadd_cout"
        }
        Opcode::IaddImm => {
            "iadd_imm"
        }
        Opcode::IaddPairwise => {
            "iadd_pairwise"
        }
        Opcode::Icmp => {
            "icmp"
        }
        Opcode::IcmpImm => {
            "icmp_imm"
        }
        Opcode::Iconcat => {
            "iconcat"
        }
        Opcode::Iconst => {
            "iconst"
        }
        Opcode::Imul => {
            "imul"
        }
        Opcode::ImulImm => {
            "imul_imm"
        }
        Opcode::Ineg => {
            "ineg"
        }
        Opcode::Insertlane => {
            "insertlane"
        }
        Opcode::Ireduce => {
            "ireduce"
        }
        Opcode::IrsubImm => {
            "irsub_imm"
        }
        Opcode::IsInvalid => {
            "is_invalid"
        }
        Opcode::IsNull => {
            "is_null"
        }
        Opcode::Ishl => {
            "ishl"
        }
        Opcode::IshlImm => {
            "ishl_imm"
        }
        Opcode::Isplit => {
            "isplit"
        }
        Opcode::Istore16 => {
            "istore16"
        }
        Opcode::Istore32 => {
            "istore32"
        }
        Opcode::Istore8 => {
            "istore8"
        }
        Opcode::Isub => {
            "isub"
        }
        Opcode::IsubBin => {
            "isub_bin"
        }
        Opcode::IsubBorrow => {
            "isub_borrow"
        }
        Opcode::IsubBout => {
            "isub_bout"
        }
        Opcode::Jump => {
            "jump"
        }
        Opcode::Load => {
            "load"
        }
        Opcode::Nearest => {
            "nearest"
        }
        Opcode::Nop => {
            "nop"
        }
        Opcode::Null => {
            "null"
        }
        Opcode::Popcnt => {
            "popcnt"
        }
        Opcode::ResumableTrap => {
            "resumable_trap"
        }
        Opcode::ResumableTrapnz => {
            "resumable_trapnz"
        }
        Opcode::Return => {
            "return"
        }
        Opcode::ReturnCall => {
            "return_call"
        }
        Opcode::ReturnCallIndirect => {
            "return_call_indirect"
        }
        Opcode::Rotl => {
            "rotl"
        }
        Opcode::RotlImm => {
            "rotl_imm"
        }
        Opcode::Rotr => {
            "rotr"
        }
        Opcode::RotrImm => {
            "rotr_imm"
        }
        Opcode::SaddSat => {
            "sadd_sat"
        }
        Opcode::ScalarToVector => {
            "scalar_to_vector"
        }
        Opcode::Sdiv => {
            "sdiv"
        }
        Opcode::SdivImm => {
            "sdiv_imm"
        }
        Opcode::Select => {
            "select"
        }
        Opcode::SelectSpectreGuard => {
            "select_spectre_guard"
        }
        Opcode::SetPinnedReg => {
            "set_pinned_reg"
        }
        Opcode::Sextend => {
            "sextend"
        }
        Opcode::Shuffle => {
            "shuffle"
        }
        Opcode::Sload16 => {
            "sload16"
        }
        Opcode::Sload16x4 => {
            "sload16x4"
        }
        Opcode::Sload32 => {
            "sload32"
        }
        Opcode::Sload32x2 => {
            "sload32x2"
        }
        Opcode::Sload8 => {
            "sload8"
        }
        Opcode::Sload8x8 => {
            "sload8x8"
        }
        Opcode::Smax => {
            "smax"
        }
        Opcode::Smin => {
            "smin"
        }
        Opcode::Smulhi => {
            "smulhi"
        }
        Opcode::Snarrow => {
            "snarrow"
        }
        Opcode::Splat => {
            "splat"
        }
        Opcode::SqmulRoundSat => {
            "sqmul_round_sat"
        }
        Opcode::Sqrt => {
            "sqrt"
        }
        Opcode::Srem => {
            "srem"
        }
        Opcode::SremImm => {
            "srem_imm"
        }
        Opcode::Sshr => {
            "sshr"
        }
        Opcode::SshrImm => {
            "sshr_imm"
        }
        Opcode::SsubSat => {
            "ssub_sat"
        }
        Opcode::StackAddr => {
            "stack_addr"
        }
        Opcode::StackLoad => {
            "stack_load"
        }
        Opcode::StackStore => {
            "stack_store"
        }
        Opcode::Store => {
            "store"
        }
        Opcode::SwidenHigh => {
            "swiden_high"
        }
        Opcode::SwidenLow => {
            "swiden_low"
        }
        Opcode::Swizzle => {
            "swizzle"
        }
        Opcode::SymbolValue => {
            "symbol_value"
        }
        Opcode::TableAddr => {
            "table_addr"
        }
        Opcode::TlsValue => {
            "tls_value"
        }
        Opcode::Trap => {
            "trap"
        }
        Opcode::Trapnz => {
            "trapnz"
        }
        Opcode::Trapz => {
            "trapz"
        }
        Opcode::Trunc => {
            "trunc"
        }
        Opcode::UaddOverflowTrap => {
            "uadd_overflow_trap"
        }
        Opcode::UaddSat => {
            "uadd_sat"
        }
        Opcode::Udiv => {
            "udiv"
        }
        Opcode::UdivImm => {
            "udiv_imm"
        }
        Opcode::Uextend => {
            "uextend"
        }
        Opcode::Uload16 => {
            "uload16"
        }
        Opcode::Uload16x4 => {
            "uload16x4"
        }
        Opcode::Uload32 => {
            "uload32"
        }
        Opcode::Uload32x2 => {
            "uload32x2"
        }
        Opcode::Uload8 => {
            "uload8"
        }
        Opcode::Uload8x8 => {
            "uload8x8"
        }
        Opcode::Umax => {
            "umax"
        }
        Opcode::Umin => {
            "umin"
        }
        Opcode::Umulhi => {
            "umulhi"
        }
        Opcode::Unarrow => {
            "unarrow"
        }
        Opcode::Urem => {
            "urem"
        }
        Opcode::UremImm => {
            "urem_imm"
        }
        Opcode::Ushr => {
            "ushr"
        }
        Opcode::UshrImm => {
            "ushr_imm"
        }
        Opcode::UsubSat => {
            "usub_sat"
        }
        Opcode::Uunarrow => {
            "uunarrow"
        }
        Opcode::UwidenHigh => {
            "uwiden_high"
        }
        Opcode::UwidenLow => {
            "uwiden_low"
        }
        Opcode::VallTrue => {
            "vall_true"
        }
        Opcode::VanyTrue => {
            "vany_true"
        }
        Opcode::Vconst => {
            "vconst"
        }
        Opcode::VhighBits => {
            "vhigh_bits"
        }
        Opcode::Vselect => {
            "vselect"
        }
    }
}

const OPCODE_HASH_TABLE: [Option<Opcode>; 256] = [
    Some(Opcode::Imul),
    Some(Opcode::TlsValue),
    None,
    Some(Opcode::Brif),
    Some(Opcode::Nearest),
    Some(Opcode::FcvtToSintSat),
    Some(Opcode::Fsub),
    Some(Opcode::Rotr),
    Some(Opcode::TableAddr),
    Some(Opcode::Iconst),
    Some(Opcode::ReturnCall),
    Some(Opcode::Umin),
    None,
    Some(Opcode::Store),
    Some(Opcode::GetFramePointer),
    Some(Opcode::UshrImm),
    Some(Opcode::Isub),
    Some(Opcode::FcvtFromSint),
    Some(Opcode::Trap),
    Some(Opcode::Sdiv),
    Some(Opcode::Srem),
    Some(Opcode::SshrImm),
    Some(Opcode::Uunarrow),
    Some(Opcode::Urem),
    Some(Opcode::Bxor),
    None,
    Some(Opcode::ExtractVector),
    Some(Opcode::Umax),
    Some(Opcode::SremImm),
    Some(Opcode::Insertlane),
    Some(Opcode::BxorNot),
    Some(Opcode::Swizzle),
    Some(Opcode::Load),
    Some(Opcode::Fadd),
    Some(Opcode::Jump),
    Some(Opcode::Null),
    Some(Opcode::Shuffle),
    Some(Opcode::Fneg),
    Some(Opcode::Umulhi),
    Some(Opcode::Ushr),
    None,
    Some(Opcode::UaddOverflowTrap),
    Some(Opcode::FcvtFromUint),
    Some(Opcode::VallTrue),
    Some(Opcode::Band),
    Some(Opcode::BxorImm),
    Some(Opcode::Fmax),
    Some(Opcode::Uload16x4),
    Some(Opcode::Ishl),
    None,
    Some(Opcode::Vconst),
    Some(Opcode::Call),
    None,
    Some(Opcode::Sqrt),
    None,
    None,
    Some(Opcode::Ceil),
    Some(Opcode::Ineg),
    Some(Opcode::FuncAddr),
    Some(Opcode::SaddSat),
    Some(Opcode::Popcnt),
    None,
    Some(Opcode::Fabs),
    Some(Opcode::Fmin),
    None,
    Some(Opcode::GlobalValue),
    Some(Opcode::Bnot),
    Some(Opcode::FmaxPseudo),
    Some(Opcode::Isplit),
    Some(Opcode::FcvtToUint),
    Some(Opcode::Trunc),
    None,
    Some(Opcode::RotlImm),
    Some(Opcode::Fcmp),
    Some(Opcode::SwidenHigh),
    Some(Opcode::Fmul),
    Some(Opcode::FcvtToSint),
    None,
    Some(Opcode::IsubBin),
    Some(Opcode::Uload8x8),
    None,
    Some(Opcode::ResumableTrapnz),
    Some(Opcode::Fdiv),
    None,
    None,
    Some(Opcode::Sextend),
    Some(Opcode::UremImm),
    Some(Opcode::AtomicLoad),
    None,
    Some(Opcode::Trapnz),
    Some(Opcode::Uload16),
    Some(Opcode::IaddImm),
    Some(Opcode::Uload32),
    Some(Opcode::IaddCarry),
    Some(Opcode::Bitrev),
    None,
    Some(Opcode::Smulhi),
    Some(Opcode::IsNull),
    None,
    None,
    None,
    None,
    None,
    Some(Opcode::BorNot),
    None,
    None,
    Some(Opcode::Sload8x8),
    None,
    None,
    None,
    None,
    Some(Opcode::SetPinnedReg),
    None,
    None,
    None,
    None,
    None,
    Some(Opcode::ImulImm),
    Some(Opcode::Ireduce),
    None,
    Some(Opcode::RotrImm),
    Some(Opcode::DynamicStackStore),
    Some(Opcode::StackStore),
    Some(Opcode::UwidenLow),
    Some(Opcode::Select),
    Some(Opcode::IaddCin),
    Some(Opcode::Istore32),
    Some(Opcode::FvpromoteLow),
    Some(Opcode::Istore16),
    None,
    Some(Opcode::Fdemote),
    Some(Opcode::BorImm),
    None,
    Some(Opcode::IcmpImm),
    Some(Opcode::Fvdemote),
    None,
    Some(Opcode::Sload16),
    Some(Opcode::Fcopysign),
    None,
    Some(Opcode::SdivImm),
    Some(Opcode::ResumableTrap),
    Some(Opcode::AvgRound),
    Some(Opcode::Sload32),
    Some(Opcode::Unarrow),
    None,
    Some(Opcode::Extractlane),
    Some(Opcode::StackAddr),
    None,
    Some(Opcode::BandImm),
    Some(Opcode::IsubBorrow),
    Some(Opcode::Return),
    None,
    Some(Opcode::Uload32x2),
    None,
    None,
    Some(Opcode::VanyTrue),
    None,
    Some(Opcode::UsubSat),
    None,
    None,
    None,
    None,
    Some(Opcode::DynamicStackLoad),
    Some(Opcode::Iconcat),
    Some(Opcode::Fence),
    None,
    None,
    None,
    None,
    Some(Opcode::Fma),
    Some(Opcode::Bitselect),
    Some(Opcode::Istore8),
    Some(Opcode::BrTable),
    Some(Opcode::F64const),
    Some(Opcode::Nop),
    Some(Opcode::StackLoad),
    Some(Opcode::IrsubImm),
    Some(Opcode::Bor),
    Some(Opcode::SqmulRoundSat),
    Some(Opcode::IsubBout),
    Some(Opcode::Debugtrap),
    Some(Opcode::Sload16x4),
    Some(Opcode::Ctz),
    Some(Opcode::IshlImm),
    Some(Opcode::BandNot),
    Some(Opcode::Bitcast),
    Some(Opcode::Cls),
    Some(Opcode::Floor),
    Some(Opcode::UwidenHigh),
    Some(Opcode::Clz),
    None,
    Some(Opcode::Uextend),
    None,
    Some(Opcode::UaddSat),
    Some(Opcode::Sload32x2),
    None,
    None,
    Some(Opcode::SelectSpectreGuard),
    Some(Opcode::Fpromote),
    None,
    None,
    None,
    Some(Opcode::SymbolValue),
    Some(Opcode::DynamicStackAddr),
    Some(Opcode::Bmask),
    Some(Opcode::GetPinnedReg),
    Some(Opcode::SsubSat),
    Some(Opcode::Vselect),
    None,
    Some(Opcode::ScalarToVector),
    Some(Opcode::AtomicRmw),
    None,
    Some(Opcode::Uload8),
    Some(Opcode::FcvtToUintSat),
    None,
    Some(Opcode::Smin),
    Some(Opcode::Trapz),
    Some(Opcode::Iabs),
    Some(Opcode::Udiv),
    None,
    Some(Opcode::AtomicCas),
    Some(Opcode::GetReturnAddress),
    None,
    Some(Opcode::SwidenLow),
    None,
    Some(Opcode::ReturnCallIndirect),
    Some(Opcode::Rotl),
    Some(Opcode::IaddPairwise),
    None,
    Some(Opcode::Smax),
    Some(Opcode::FminPseudo),
    None,
    Some(Opcode::F32const),
    Some(Opcode::UdivImm),
    None,
    Some(Opcode::Splat),
    Some(Opcode::IaddCout),
    Some(Opcode::Snarrow),
    Some(Opcode::CallIndirect),
    Some(Opcode::Sload8),
    None,
    Some(Opcode::VhighBits),
    Some(Opcode::IsInvalid),
    None,
    Some(Opcode::Iadd),
    Some(Opcode::Icmp),
    None,
    Some(Opcode::FcvtLowFromSint),
    None,
    None,
    Some(Opcode::GetStackPointer),
    Some(Opcode::Bswap),
    None,
    Some(Opcode::Sshr),
    Some(Opcode::AtomicStore),
    None,
];


// Table of opcode constraints.
const OPCODE_CONSTRAINTS: [OpcodeConstraints; 181] = [
    // Jump: fixed_results=0, use_typevar_operand=false, requires_typevar_operand=false, fixed_values=0
    // Constraints=[]
    OpcodeConstraints {
        flags: 0x00,
        typeset_offset: 255,
        constraint_offset: 0,
    },
    // Brif: fixed_results=0, use_typevar_operand=true, requires_typevar_operand=true, fixed_values=1
    // Constraints=['Same']
    // Polymorphic over TypeSet(lanes={1}, ints={8, 16, 32, 64, 128})
    OpcodeConstraints {
        flags: 0x38,
        typeset_offset: 0,
        constraint_offset: 0,
    },
    // BrTable: fixed_results=0, use_typevar_operand=false, requires_typevar_operand=false, fixed_values=1
    // Constraints=['Concrete(ir::types::I32)']
    OpcodeConstraints {
        flags: 0x20,
        typeset_offset: 255,
        constraint_offset: 3,
    },
    // Debugtrap: fixed_results=0, use_typevar_operand=false, requires_typevar_operand=false, fixed_values=0
    // Constraints=[]
    OpcodeConstraints {
        flags: 0x00,
        typeset_offset: 255,
        constraint_offset: 0,
    },
    // Trap: fixed_results=0, use_typevar_operand=false, requires_typevar_operand=false, fixed_values=0
    // Constraints=[]
    OpcodeConstraints {
        flags: 0x00,
        typeset_offset: 255,
        constraint_offset: 0,
    },
    // Trapz: fixed_results=0, use_typevar_operand=true, requires_typevar_operand=true, fixed_values=1
    // Constraints=['Same']
    // Polymorphic over TypeSet(lanes={1}, ints={8, 16, 32, 64, 128})
    OpcodeConstraints {
        flags: 0x38,
        typeset_offset: 0,
        constraint_offset: 0,
    },
    // ResumableTrap: fixed_results=0, use_typevar_operand=false, requires_typevar_operand=false, fixed_values=0
    // Constraints=[]
    OpcodeConstraints {
        flags: 0x00,
        typeset_offset: 255,
        constraint_offset: 0,
    },
    // Trapnz: fixed_results=0, use_typevar_operand=true, requires_typevar_operand=true, fixed_values=1
    // Constraints=['Same']
    // Polymorphic over TypeSet(lanes={1}, ints={8, 16, 32, 64, 128})
    OpcodeConstraints {
        flags: 0x38,
        typeset_offset: 0,
        constraint_offset: 0,
    },
    // ResumableTrapnz: fixed_results=0, use_typevar_operand=true, requires_typevar_operand=true, fixed_values=1
    // Constraints=['Same']
    // Polymorphic over TypeSet(lanes={1}, ints={8, 16, 32, 64, 128})
    OpcodeConstraints {
        flags: 0x38,
        typeset_offset: 0,
        constraint_offset: 0,
    },
    // Return: fixed_results=0, use_typevar_operand=false, requires_typevar_operand=false, fixed_values=0
    // Constraints=[]
    OpcodeConstraints {
        flags: 0x00,
        typeset_offset: 255,
        constraint_offset: 0,
    },
    // Call: fixed_results=0, use_typevar_operand=false, requires_typevar_operand=false, fixed_values=0
    // Constraints=[]
    OpcodeConstraints {
        flags: 0x00,
        typeset_offset: 255,
        constraint_offset: 0,
    },
    // CallIndirect: fixed_results=0, use_typevar_operand=true, requires_typevar_operand=true, fixed_values=1
    // Constraints=['Same']
    // Polymorphic over TypeSet(lanes={1}, ints={32, 64}, refs={32, 64})
    OpcodeConstraints {
        flags: 0x38,
        typeset_offset: 1,
        constraint_offset: 0,
    },
    // ReturnCall: fixed_results=0, use_typevar_operand=false, requires_typevar_operand=false, fixed_values=0
    // Constraints=[]
    OpcodeConstraints {
        flags: 0x00,
        typeset_offset: 255,
        constraint_offset: 0,
    },
    // ReturnCallIndirect: fixed_results=0, use_typevar_operand=true, requires_typevar_operand=true, fixed_values=1
    // Constraints=['Same']
    // Polymorphic over TypeSet(lanes={1}, ints={32, 64}, refs={32, 64})
    OpcodeConstraints {
        flags: 0x38,
        typeset_offset: 1,
        constraint_offset: 0,
    },
    // FuncAddr: fixed_results=1, use_typevar_operand=false, requires_typevar_operand=false, fixed_values=0
    // Constraints=['Same']
    // Polymorphic over TypeSet(lanes={1}, ints={32, 64}, refs={32, 64})
    OpcodeConstraints {
        flags: 0x01,
        typeset_offset: 1,
        constraint_offset: 0,
    },
    // Splat: fixed_results=1, use_typevar_operand=false, requires_typevar_operand=false, fixed_values=1
    // Constraints=['Same', 'LaneOf']
    // Polymorphic over TypeSet(lanes={2, 4, 8, 16, 32, 64, 128, 256}, ints={8, 16, 32, 64, 128}, floats={32, 64})
    OpcodeConstraints {
        flags: 0x21,
        typeset_offset: 2,
        constraint_offset: 4,
    },
    // Swizzle: fixed_results=1, use_typevar_operand=false, requires_typevar_operand=false, fixed_values=2
    // Constraints=['Concrete(ir::types::I8X16)', 'Concrete(ir::types::I8X16)', 'Concrete(ir::types::I8X16)']
    OpcodeConstraints {
        flags: 0x41,
        typeset_offset: 255,
        constraint_offset: 6,
    },
    // Insertlane: fixed_results=1, use_typevar_operand=true, requires_typevar_operand=false, fixed_values=2
    // Constraints=['Same', 'Same', 'LaneOf']
    // Polymorphic over TypeSet(lanes={2, 4, 8, 16, 32, 64, 128, 256}, ints={8, 16, 32, 64, 128}, floats={32, 64})
    OpcodeConstraints {
        flags: 0x49,
        typeset_offset: 2,
        constraint_offset: 9,
    },
    // Extractlane: fixed_results=1, use_typevar_operand=true, requires_typevar_operand=true, fixed_values=1
    // Constraints=['LaneOf', 'Same']
    // Polymorphic over TypeSet(lanes={2, 4, 8, 16, 32, 64, 128, 256}, ints={8, 16, 32, 64, 128}, floats={32, 64})
    OpcodeConstraints {
        flags: 0x39,
        typeset_offset: 2,
        constraint_offset: 11,
    },
    // Smin: fixed_results=1, use_typevar_operand=true, requires_typevar_operand=false, fixed_values=2
    // Constraints=['Same', 'Same', 'Same']
    // Polymorphic over TypeSet(lanes={1, 2, 4, 8, 16, 32, 64, 128, 256}, ints={8, 16, 32, 64, 128})
    OpcodeConstraints {
        flags: 0x49,
        typeset_offset: 3,
        constraint_offset: 0,
    },
    // Umin: fixed_results=1, use_typevar_operand=true, requires_typevar_operand=false, fixed_values=2
    // Constraints=['Same', 'Same', 'Same']
    // Polymorphic over TypeSet(lanes={1, 2, 4, 8, 16, 32, 64, 128, 256}, ints={8, 16, 32, 64, 128})
    OpcodeConstraints {
        flags: 0x49,
        typeset_offset: 3,
        constraint_offset: 0,
    },
    // Smax: fixed_results=1, use_typevar_operand=true, requires_typevar_operand=false, fixed_values=2
    // Constraints=['Same', 'Same', 'Same']
    // Polymorphic over TypeSet(lanes={1, 2, 4, 8, 16, 32, 64, 128, 256}, ints={8, 16, 32, 64, 128})
    OpcodeConstraints {
        flags: 0x49,
        typeset_offset: 3,
        constraint_offset: 0,
    },
    // Umax: fixed_results=1, use_typevar_operand=true, requires_typevar_operand=false, fixed_values=2
    // Constraints=['Same', 'Same', 'Same']
    // Polymorphic over TypeSet(lanes={1, 2, 4, 8, 16, 32, 64, 128, 256}, ints={8, 16, 32, 64, 128})
    OpcodeConstraints {
        flags: 0x49,
        typeset_offset: 3,
        constraint_offset: 0,
    },
    // AvgRound: fixed_results=1, use_typevar_operand=true, requires_typevar_operand=false, fixed_values=2
    // Constraints=['Same', 'Same', 'Same']
    // Polymorphic over TypeSet(lanes={2, 4, 8, 16, 32, 64, 128, 256}, ints={8, 16, 32, 64, 128})
    OpcodeConstraints {
        flags: 0x49,
        typeset_offset: 4,
        constraint_offset: 0,
    },
    // UaddSat: fixed_results=1, use_typevar_operand=true, requires_typevar_operand=false, fixed_values=2
    // Constraints=['Same', 'Same', 'Same']
    // Polymorphic over TypeSet(lanes={2, 4, 8, 16, 32, 64, 128, 256}, ints={8, 16, 32, 64, 128})
    OpcodeConstraints {
        flags: 0x49,
        typeset_offset: 4,
        constraint_offset: 0,
    },
    // SaddSat: fixed_results=1, use_typevar_operand=true, requires_typevar_operand=false, fixed_values=2
    // Constraints=['Same', 'Same', 'Same']
    // Polymorphic over TypeSet(lanes={2, 4, 8, 16, 32, 64, 128, 256}, ints={8, 16, 32, 64, 128})
    OpcodeConstraints {
        flags: 0x49,
        typeset_offset: 4,
        constraint_offset: 0,
    },
    // UsubSat: fixed_results=1, use_typevar_operand=true, requires_typevar_operand=false, fixed_values=2
    // Constraints=['Same', 'Same', 'Same']
    // Polymorphic over TypeSet(lanes={2, 4, 8, 16, 32, 64, 128, 256}, ints={8, 16, 32, 64, 128})
    OpcodeConstraints {
        flags: 0x49,
        typeset_offset: 4,
        constraint_offset: 0,
    },
    // SsubSat: fixed_results=1, use_typevar_operand=true, requires_typevar_operand=false, fixed_values=2
    // Constraints=['Same', 'Same', 'Same']
    // Polymorphic over TypeSet(lanes={2, 4, 8, 16, 32, 64, 128, 256}, ints={8, 16, 32, 64, 128})
    OpcodeConstraints {
        flags: 0x49,
        typeset_offset: 4,
        constraint_offset: 0,
    },
    // Load: fixed_results=1, use_typevar_operand=false, requires_typevar_operand=false, fixed_values=1
    // Constraints=['Same', 'Free(1)']
    // Polymorphic over TypeSet(lanes={1, 2, 4, 8, 16, 32, 64, 128, 256}, ints={8, 16, 32, 64, 128}, floats={32, 64}, refs={32, 64})
    OpcodeConstraints {
        flags: 0x21,
        typeset_offset: 5,
        constraint_offset: 12,
    },
    // Store: fixed_results=0, use_typevar_operand=true, requires_typevar_operand=true, fixed_values=2
    // Constraints=['Same', 'Free(1)']
    // Polymorphic over TypeSet(lanes={1, 2, 4, 8, 16, 32, 64, 128, 256}, ints={8, 16, 32, 64, 128}, floats={32, 64}, refs={32, 64})
    OpcodeConstraints {
        flags: 0x58,
        typeset_offset: 5,
        constraint_offset: 12,
    },
    // Uload8: fixed_results=1, use_typevar_operand=false, requires_typevar_operand=false, fixed_values=1
    // Constraints=['Same', 'Free(1)']
    // Polymorphic over TypeSet(lanes={1}, ints={16, 32, 64})
    OpcodeConstraints {
        flags: 0x21,
        typeset_offset: 6,
        constraint_offset: 12,
    },
    // Sload8: fixed_results=1, use_typevar_operand=false, requires_typevar_operand=false, fixed_values=1
    // Constraints=['Same', 'Free(1)']
    // Polymorphic over TypeSet(lanes={1}, ints={16, 32, 64})
    OpcodeConstraints {
        flags: 0x21,
        typeset_offset: 6,
        constraint_offset: 12,
    },
    // Istore8: fixed_results=0, use_typevar_operand=true, requires_typevar_operand=true, fixed_values=2
    // Constraints=['Same', 'Free(1)']
    // Polymorphic over TypeSet(lanes={1}, ints={16, 32, 64})
    OpcodeConstraints {
        flags: 0x58,
        typeset_offset: 6,
        constraint_offset: 12,
    },
    // Uload16: fixed_results=1, use_typevar_operand=false, requires_typevar_operand=false, fixed_values=1
    // Constraints=['Same', 'Free(1)']
    // Polymorphic over TypeSet(lanes={1}, ints={32, 64})
    OpcodeConstraints {
        flags: 0x21,
        typeset_offset: 7,
        constraint_offset: 12,
    },
    // Sload16: fixed_results=1, use_typevar_operand=false, requires_typevar_operand=false, fixed_values=1
    // Constraints=['Same', 'Free(1)']
    // Polymorphic over TypeSet(lanes={1}, ints={32, 64})
    OpcodeConstraints {
        flags: 0x21,
        typeset_offset: 7,
        constraint_offset: 12,
    },
    // Istore16: fixed_results=0, use_typevar_operand=true, requires_typevar_operand=true, fixed_values=2
    // Constraints=['Same', 'Free(1)']
    // Polymorphic over TypeSet(lanes={1}, ints={32, 64})
    OpcodeConstraints {
        flags: 0x58,
        typeset_offset: 7,
        constraint_offset: 12,
    },
    // Uload32: fixed_results=1, use_typevar_operand=true, requires_typevar_operand=true, fixed_values=1
    // Constraints=['Concrete(ir::types::I64)', 'Same']
    // Polymorphic over TypeSet(lanes={1}, ints={32, 64}, refs={32, 64})
    OpcodeConstraints {
        flags: 0x39,
        typeset_offset: 1,
        constraint_offset: 14,
    },
    // Sload32: fixed_results=1, use_typevar_operand=true, requires_typevar_operand=true, fixed_values=1
    // Constraints=['Concrete(ir::types::I64)', 'Same']
    // Polymorphic over TypeSet(lanes={1}, ints={32, 64}, refs={32, 64})
    OpcodeConstraints {
        flags: 0x39,
        typeset_offset: 1,
        constraint_offset: 14,
    },
    // Istore32: fixed_results=0, use_typevar_operand=true, requires_typevar_operand=true, fixed_values=2
    // Constraints=['Concrete(ir::types::I64)', 'Free(1)']
    // Polymorphic over TypeSet(lanes={1}, ints={64})
    OpcodeConstraints {
        flags: 0x58,
        typeset_offset: 8,
        constraint_offset: 16,
    },
    // Uload8x8: fixed_results=1, use_typevar_operand=true, requires_typevar_operand=true, fixed_values=1
    // Constraints=['Concrete(ir::types::I16X8)', 'Same']
    // Polymorphic over TypeSet(lanes={1}, ints={32, 64}, refs={32, 64})
    OpcodeConstraints {
        flags: 0x39,
        typeset_offset: 1,
        constraint_offset: 18,
    },
    // Sload8x8: fixed_results=1, use_typevar_operand=true, requires_typevar_operand=true, fixed_values=1
    // Constraints=['Concrete(ir::types::I16X8)', 'Same']
    // Polymorphic over TypeSet(lanes={1}, ints={32, 64}, refs={32, 64})
    OpcodeConstraints {
        flags: 0x39,
        typeset_offset: 1,
        constraint_offset: 18,
    },
    // Uload16x4: fixed_results=1, use_typevar_operand=true, requires_typevar_operand=true, fixed_values=1
    // Constraints=['Concrete(ir::types::I32X4)', 'Same']
    // Polymorphic over TypeSet(lanes={1}, ints={32, 64}, refs={32, 64})
    OpcodeConstraints {
        flags: 0x39,
        typeset_offset: 1,
        constraint_offset: 20,
    },
    // Sload16x4: fixed_results=1, use_typevar_operand=true, requires_typevar_operand=true, fixed_values=1
    // Constraints=['Concrete(ir::types::I32X4)', 'Same']
    // Polymorphic over TypeSet(lanes={1}, ints={32, 64}, refs={32, 64})
    OpcodeConstraints {
        flags: 0x39,
        typeset_offset: 1,
        constraint_offset: 20,
    },
    // Uload32x2: fixed_results=1, use_typevar_operand=true, requires_typevar_operand=true, fixed_values=1
    // Constraints=['Concrete(ir::types::I64X2)', 'Same']
    // Polymorphic over TypeSet(lanes={1}, ints={32, 64}, refs={32, 64})
    OpcodeConstraints {
        flags: 0x39,
        typeset_offset: 1,
        constraint_offset: 22,
    },
    // Sload32x2: fixed_results=1, use_typevar_operand=true, requires_typevar_operand=true, fixed_values=1
    // Constraints=['Concrete(ir::types::I64X2)', 'Same']
    // Polymorphic over TypeSet(lanes={1}, ints={32, 64}, refs={32, 64})
    OpcodeConstraints {
        flags: 0x39,
        typeset_offset: 1,
        constraint_offset: 22,
    },
    // StackLoad: fixed_results=1, use_typevar_operand=false, requires_typevar_operand=false, fixed_values=0
    // Constraints=['Same']
    // Polymorphic over TypeSet(lanes={1, 2, 4, 8, 16, 32, 64, 128, 256}, ints={8, 16, 32, 64, 128}, floats={32, 64}, refs={32, 64})
    OpcodeConstraints {
        flags: 0x01,
        typeset_offset: 5,
        constraint_offset: 0,
    },
    // StackStore: fixed_results=0, use_typevar_operand=true, requires_typevar_operand=true, fixed_values=1
    // Constraints=['Same']
    // Polymorphic over TypeSet(lanes={1, 2, 4, 8, 16, 32, 64, 128, 256}, ints={8, 16, 32, 64, 128}, floats={32, 64}, refs={32, 64})
    OpcodeConstraints {
        flags: 0x38,
        typeset_offset: 5,
        constraint_offset: 0,
    },
    // StackAddr: fixed_results=1, use_typevar_operand=false, requires_typevar_operand=false, fixed_values=0
    // Constraints=['Same']
    // Polymorphic over TypeSet(lanes={1}, ints={32, 64}, refs={32, 64})
    OpcodeConstraints {
        flags: 0x01,
        typeset_offset: 1,
        constraint_offset: 0,
    },
    // DynamicStackLoad: fixed_results=1, use_typevar_operand=false, requires_typevar_operand=false, fixed_values=0
    // Constraints=['Same']
    // Polymorphic over TypeSet(lanes={1, 2, 4, 8, 16, 32, 64, 128, 256}, ints={8, 16, 32, 64, 128}, floats={32, 64}, refs={32, 64})
    OpcodeConstraints {
        flags: 0x01,
        typeset_offset: 5,
        constraint_offset: 0,
    },
    // DynamicStackStore: fixed_results=0, use_typevar_operand=true, requires_typevar_operand=true, fixed_values=1
    // Constraints=['Same']
    // Polymorphic over TypeSet(lanes={1, 2, 4, 8, 16, 32, 64, 128, 256}, ints={8, 16, 32, 64, 128}, floats={32, 64}, refs={32, 64})
    OpcodeConstraints {
        flags: 0x38,
        typeset_offset: 5,
        constraint_offset: 0,
    },
    // DynamicStackAddr: fixed_results=1, use_typevar_operand=false, requires_typevar_operand=false, fixed_values=0
    // Constraints=['Same']
    // Polymorphic over TypeSet(lanes={1}, ints={32, 64}, refs={32, 64})
    OpcodeConstraints {
        flags: 0x01,
        typeset_offset: 1,
        constraint_offset: 0,
    },
    // GlobalValue: fixed_results=1, use_typevar_operand=false, requires_typevar_operand=false, fixed_values=0
    // Constraints=['Same']
    // Polymorphic over TypeSet(lanes={1, 2, 4, 8, 16, 32, 64, 128, 256}, ints={8, 16, 32, 64, 128}, floats={32, 64}, refs={32, 64})
    OpcodeConstraints {
        flags: 0x01,
        typeset_offset: 5,
        constraint_offset: 0,
    },
    // SymbolValue: fixed_results=1, use_typevar_operand=false, requires_typevar_operand=false, fixed_values=0
    // Constraints=['Same']
    // Polymorphic over TypeSet(lanes={1, 2, 4, 8, 16, 32, 64, 128, 256}, ints={8, 16, 32, 64, 128}, floats={32, 64}, refs={32, 64})
    OpcodeConstraints {
        flags: 0x01,
        typeset_offset: 5,
        constraint_offset: 0,
    },
    // TlsValue: fixed_results=1, use_typevar_operand=false, requires_typevar_operand=false, fixed_values=0
    // Constraints=['Same']
    // Polymorphic over TypeSet(lanes={1, 2, 4, 8, 16, 32, 64, 128, 256}, ints={8, 16, 32, 64, 128}, floats={32, 64}, refs={32, 64})
    OpcodeConstraints {
        flags: 0x01,
        typeset_offset: 5,
        constraint_offset: 0,
    },
    // GetPinnedReg: fixed_results=1, use_typevar_operand=false, requires_typevar_operand=false, fixed_values=0
    // Constraints=['Same']
    // Polymorphic over TypeSet(lanes={1}, ints={32, 64}, refs={32, 64})
    OpcodeConstraints {
        flags: 0x01,
        typeset_offset: 1,
        constraint_offset: 0,
    },
    // SetPinnedReg: fixed_results=0, use_typevar_operand=true, requires_typevar_operand=true, fixed_values=1
    // Constraints=['Same']
    // Polymorphic over TypeSet(lanes={1}, ints={32, 64}, refs={32, 64})
    OpcodeConstraints {
        flags: 0x38,
        typeset_offset: 1,
        constraint_offset: 0,
    },
    // GetFramePointer: fixed_results=1, use_typevar_operand=false, requires_typevar_operand=false, fixed_values=0
    // Constraints=['Same']
    // Polymorphic over TypeSet(lanes={1}, ints={32, 64}, refs={32, 64})
    OpcodeConstraints {
        flags: 0x01,
        typeset_offset: 1,
        constraint_offset: 0,
    },
    // GetStackPointer: fixed_results=1, use_typevar_operand=false, requires_typevar_operand=false, fixed_values=0
    // Constraints=['Same']
    // Polymorphic over TypeSet(lanes={1}, ints={32, 64}, refs={32, 64})
    OpcodeConstraints {
        flags: 0x01,
        typeset_offset: 1,
        constraint_offset: 0,
    },
    // GetReturnAddress: fixed_results=1, use_typevar_operand=false, requires_typevar_operand=false, fixed_values=0
    // Constraints=['Same']
    // Polymorphic over TypeSet(lanes={1}, ints={32, 64}, refs={32, 64})
    OpcodeConstraints {
        flags: 0x01,
        typeset_offset: 1,
        constraint_offset: 0,
    },
    // TableAddr: fixed_results=1, use_typevar_operand=false, requires_typevar_operand=false, fixed_values=1
    // Constraints=['Same', 'Free(7)']
    // Polymorphic over TypeSet(lanes={1}, ints={32, 64}, refs={32, 64})
    OpcodeConstraints {
        flags: 0x21,
        typeset_offset: 1,
        constraint_offset: 23,
    },
    // Iconst: fixed_results=1, use_typevar_operand=false, requires_typevar_operand=false, fixed_values=0
    // Constraints=['Same']
    // Polymorphic over TypeSet(lanes={1, 2, 4, 8, 16, 32, 64, 128, 256}, ints={8, 16, 32, 64})
    OpcodeConstraints {
        flags: 0x01,
        typeset_offset: 9,
        constraint_offset: 0,
    },
    // F32const: fixed_results=1, use_typevar_operand=false, requires_typevar_operand=false, fixed_values=0
    // Constraints=['Concrete(ir::types::F32)']
    OpcodeConstraints {
        flags: 0x01,
        typeset_offset: 255,
        constraint_offset: 25,
    },
    // F64const: fixed_results=1, use_typevar_operand=false, requires_typevar_operand=false, fixed_values=0
    // Constraints=['Concrete(ir::types::F64)']
    OpcodeConstraints {
        flags: 0x01,
        typeset_offset: 255,
        constraint_offset: 26,
    },
    // Vconst: fixed_results=1, use_typevar_operand=false, requires_typevar_operand=false, fixed_values=0
    // Constraints=['Same']
    // Polymorphic over TypeSet(lanes={2, 4, 8, 16, 32, 64, 128, 256}, ints={8, 16, 32, 64, 128}, floats={32, 64})
    OpcodeConstraints {
        flags: 0x01,
        typeset_offset: 10,
        constraint_offset: 0,
    },
    // Shuffle: fixed_results=1, use_typevar_operand=false, requires_typevar_operand=false, fixed_values=2
    // Constraints=['Concrete(ir::types::I8X16)', 'Concrete(ir::types::I8X16)', 'Concrete(ir::types::I8X16)']
    OpcodeConstraints {
        flags: 0x41,
        typeset_offset: 255,
        constraint_offset: 6,
    },
    // Null: fixed_results=1, use_typevar_operand=false, requires_typevar_operand=false, fixed_values=0
    // Constraints=['Same']
    // Polymorphic over TypeSet(lanes={1}, refs={32, 64})
    OpcodeConstraints {
        flags: 0x01,
        typeset_offset: 11,
        constraint_offset: 0,
    },
    // Nop: fixed_results=0, use_typevar_operand=false, requires_typevar_operand=false, fixed_values=0
    // Constraints=[]
    OpcodeConstraints {
        flags: 0x00,
        typeset_offset: 255,
        constraint_offset: 0,
    },
    // Select: fixed_results=1, use_typevar_operand=true, requires_typevar_operand=false, fixed_values=3
    // Constraints=['Same', 'Free(0)', 'Same', 'Same']
    // Polymorphic over TypeSet(lanes={1, 2, 4, 8, 16, 32, 64, 128, 256}, ints={8, 16, 32, 64, 128}, floats={32, 64}, refs={32, 64})
    OpcodeConstraints {
        flags: 0x69,
        typeset_offset: 12,
        constraint_offset: 27,
    },
    // SelectSpectreGuard: fixed_results=1, use_typevar_operand=true, requires_typevar_operand=false, fixed_values=3
    // Constraints=['Same', 'Free(0)', 'Same', 'Same']
    // Polymorphic over TypeSet(lanes={1, 2, 4, 8, 16, 32, 64, 128, 256}, ints={8, 16, 32, 64, 128}, floats={32, 64}, refs={32, 64})
    OpcodeConstraints {
        flags: 0x69,
        typeset_offset: 12,
        constraint_offset: 27,
    },
    // Bitselect: fixed_results=1, use_typevar_operand=true, requires_typevar_operand=false, fixed_values=3
    // Constraints=['Same', 'Same', 'Same', 'Same']
    // Polymorphic over TypeSet(lanes={1, 2, 4, 8, 16, 32, 64, 128, 256}, ints={8, 16, 32, 64, 128}, floats={32, 64}, refs={32, 64})
    OpcodeConstraints {
        flags: 0x69,
        typeset_offset: 12,
        constraint_offset: 29,
    },
    // Vselect: fixed_results=1, use_typevar_operand=true, requires_typevar_operand=false, fixed_values=3
    // Constraints=['Same', 'AsBool', 'Same', 'Same']
    // Polymorphic over TypeSet(lanes={2, 4, 8, 16, 32, 64, 128, 256}, ints={8, 16, 32, 64, 128}, floats={32, 64})
    OpcodeConstraints {
        flags: 0x69,
        typeset_offset: 10,
        constraint_offset: 32,
    },
    // VanyTrue: fixed_results=1, use_typevar_operand=true, requires_typevar_operand=true, fixed_values=1
    // Constraints=['Concrete(ir::types::I8)', 'Same']
    // Polymorphic over TypeSet(lanes={2, 4, 8, 16, 32, 64, 128, 256}, ints={8, 16, 32, 64, 128}, floats={32, 64})
    OpcodeConstraints {
        flags: 0x39,
        typeset_offset: 10,
        constraint_offset: 36,
    },
    // VallTrue: fixed_results=1, use_typevar_operand=true, requires_typevar_operand=true, fixed_values=1
    // Constraints=['Concrete(ir::types::I8)', 'Same']
    // Polymorphic over TypeSet(lanes={2, 4, 8, 16, 32, 64, 128, 256}, ints={8, 16, 32, 64, 128}, floats={32, 64})
    OpcodeConstraints {
        flags: 0x39,
        typeset_offset: 10,
        constraint_offset: 36,
    },
    // VhighBits: fixed_results=1, use_typevar_operand=false, requires_typevar_operand=false, fixed_values=1
    // Constraints=['Same', 'Free(10)']
    // Polymorphic over TypeSet(lanes={1, 2, 4, 8, 16, 32, 64, 128, 256}, ints={8, 16, 32, 64, 128})
    OpcodeConstraints {
        flags: 0x21,
        typeset_offset: 13,
        constraint_offset: 37,
    },
    // Icmp: fixed_results=1, use_typevar_operand=true, requires_typevar_operand=true, fixed_values=2
    // Constraints=['AsBool', 'Same', 'Same']
    // Polymorphic over TypeSet(lanes={1, 2, 4, 8, 16, 32, 64, 128, 256}, ints={8, 16, 32, 64, 128})
    OpcodeConstraints {
        flags: 0x59,
        typeset_offset: 13,
        constraint_offset: 33,
    },
    // IcmpImm: fixed_results=1, use_typevar_operand=true, requires_typevar_operand=true, fixed_values=1
    // Constraints=['Concrete(ir::types::I8)', 'Same']
    // Polymorphic over TypeSet(lanes={1}, ints={8, 16, 32, 64, 128})
    OpcodeConstraints {
        flags: 0x39,
        typeset_offset: 0,
        constraint_offset: 36,
    },
    // Iadd: fixed_results=1, use_typevar_operand=true, requires_typevar_operand=false, fixed_values=2
    // Constraints=['Same', 'Same', 'Same']
    // Polymorphic over TypeSet(lanes={1, 2, 4, 8, 16, 32, 64, 128, 256}, ints={8, 16, 32, 64, 128})
    OpcodeConstraints {
        flags: 0x49,
        typeset_offset: 13,
        constraint_offset: 0,
    },
    // Isub: fixed_results=1, use_typevar_operand=true, requires_typevar_operand=false, fixed_values=2
    // Constraints=['Same', 'Same', 'Same']
    // Polymorphic over TypeSet(lanes={1, 2, 4, 8, 16, 32, 64, 128, 256}, ints={8, 16, 32, 64, 128})
    OpcodeConstraints {
        flags: 0x49,
        typeset_offset: 13,
        constraint_offset: 0,
    },
    // Ineg: fixed_results=1, use_typevar_operand=true, requires_typevar_operand=false, fixed_values=1
    // Constraints=['Same', 'Same']
    // Polymorphic over TypeSet(lanes={1, 2, 4, 8, 16, 32, 64, 128, 256}, ints={8, 16, 32, 64, 128})
    OpcodeConstraints {
        flags: 0x29,
        typeset_offset: 13,
        constraint_offset: 0,
    },
    // Iabs: fixed_results=1, use_typevar_operand=true, requires_typevar_operand=false, fixed_values=1
    // Constraints=['Same', 'Same']
    // Polymorphic over TypeSet(lanes={1, 2, 4, 8, 16, 32, 64, 128, 256}, ints={8, 16, 32, 64, 128})
    OpcodeConstraints {
        flags: 0x29,
        typeset_offset: 13,
        constraint_offset: 0,
    },
    // Imul: fixed_results=1, use_typevar_operand=true, requires_typevar_operand=false, fixed_values=2
    // Constraints=['Same', 'Same', 'Same']
    // Polymorphic over TypeSet(lanes={1, 2, 4, 8, 16, 32, 64, 128, 256}, ints={8, 16, 32, 64, 128})
    OpcodeConstraints {
        flags: 0x49,
        typeset_offset: 13,
        constraint_offset: 0,
    },
    // Umulhi: fixed_results=1, use_typevar_operand=true, requires_typevar_operand=false, fixed_values=2
    // Constraints=['Same', 'Same', 'Same']
    // Polymorphic over TypeSet(lanes={1, 2, 4, 8, 16, 32, 64, 128, 256}, ints={8, 16, 32, 64, 128})
    OpcodeConstraints {
        flags: 0x49,
        typeset_offset: 13,
        constraint_offset: 0,
    },
    // Smulhi: fixed_results=1, use_typevar_operand=true, requires_typevar_operand=false, fixed_values=2
    // Constraints=['Same', 'Same', 'Same']
    // Polymorphic over TypeSet(lanes={1, 2, 4, 8, 16, 32, 64, 128, 256}, ints={8, 16, 32, 64, 128})
    OpcodeConstraints {
        flags: 0x49,
        typeset_offset: 13,
        constraint_offset: 0,
    },
    // SqmulRoundSat: fixed_results=1, use_typevar_operand=true, requires_typevar_operand=false, fixed_values=2
    // Constraints=['Same', 'Same', 'Same']
    // Polymorphic over TypeSet(lanes={4, 8}, ints={16, 32})
    OpcodeConstraints {
        flags: 0x49,
        typeset_offset: 14,
        constraint_offset: 0,
    },
    // Udiv: fixed_results=1, use_typevar_operand=true, requires_typevar_operand=false, fixed_values=2
    // Constraints=['Same', 'Same', 'Same']
    // Polymorphic over TypeSet(lanes={1}, ints={8, 16, 32, 64, 128})
    OpcodeConstraints {
        flags: 0x49,
        typeset_offset: 0,
        constraint_offset: 0,
    },
    // Sdiv: fixed_results=1, use_typevar_operand=true, requires_typevar_operand=false, fixed_values=2
    // Constraints=['Same', 'Same', 'Same']
    // Polymorphic over TypeSet(lanes={1}, ints={8, 16, 32, 64, 128})
    OpcodeConstraints {
        flags: 0x49,
        typeset_offset: 0,
        constraint_offset: 0,
    },
    // Urem: fixed_results=1, use_typevar_operand=true, requires_typevar_operand=false, fixed_values=2
    // Constraints=['Same', 'Same', 'Same']
    // Polymorphic over TypeSet(lanes={1}, ints={8, 16, 32, 64, 128})
    OpcodeConstraints {
        flags: 0x49,
        typeset_offset: 0,
        constraint_offset: 0,
    },
    // Srem: fixed_results=1, use_typevar_operand=true, requires_typevar_operand=false, fixed_values=2
    // Constraints=['Same', 'Same', 'Same']
    // Polymorphic over TypeSet(lanes={1}, ints={8, 16, 32, 64, 128})
    OpcodeConstraints {
        flags: 0x49,
        typeset_offset: 0,
        constraint_offset: 0,
    },
    // IaddImm: fixed_results=1, use_typevar_operand=true, requires_typevar_operand=false, fixed_values=1
    // Constraints=['Same', 'Same']
    // Polymorphic over TypeSet(lanes={1}, ints={8, 16, 32, 64, 128})
    OpcodeConstraints {
        flags: 0x29,
        typeset_offset: 0,
        constraint_offset: 0,
    },
    // ImulImm: fixed_results=1, use_typevar_operand=true, requires_typevar_operand=false, fixed_values=1
    // Constraints=['Same', 'Same']
    // Polymorphic over TypeSet(lanes={1}, ints={8, 16, 32, 64, 128})
    OpcodeConstraints {
        flags: 0x29,
        typeset_offset: 0,
        constraint_offset: 0,
    },
    // UdivImm: fixed_results=1, use_typevar_operand=true, requires_typevar_operand=false, fixed_values=1
    // Constraints=['Same', 'Same']
    // Polymorphic over TypeSet(lanes={1}, ints={8, 16, 32, 64, 128})
    OpcodeConstraints {
        flags: 0x29,
        typeset_offset: 0,
        constraint_offset: 0,
    },
    // SdivImm: fixed_results=1, use_typevar_operand=true, requires_typevar_operand=false, fixed_values=1
    // Constraints=['Same', 'Same']
    // Polymorphic over TypeSet(lanes={1}, ints={8, 16, 32, 64, 128})
    OpcodeConstraints {
        flags: 0x29,
        typeset_offset: 0,
        constraint_offset: 0,
    },
    // UremImm: fixed_results=1, use_typevar_operand=true, requires_typevar_operand=false, fixed_values=1
    // Constraints=['Same', 'Same']
    // Polymorphic over TypeSet(lanes={1}, ints={8, 16, 32, 64, 128})
    OpcodeConstraints {
        flags: 0x29,
        typeset_offset: 0,
        constraint_offset: 0,
    },
    // SremImm: fixed_results=1, use_typevar_operand=true, requires_typevar_operand=false, fixed_values=1
    // Constraints=['Same', 'Same']
    // Polymorphic over TypeSet(lanes={1}, ints={8, 16, 32, 64, 128})
    OpcodeConstraints {
        flags: 0x29,
        typeset_offset: 0,
        constraint_offset: 0,
    },
    // IrsubImm: fixed_results=1, use_typevar_operand=true, requires_typevar_operand=false, fixed_values=1
    // Constraints=['Same', 'Same']
    // Polymorphic over TypeSet(lanes={1}, ints={8, 16, 32, 64, 128})
    OpcodeConstraints {
        flags: 0x29,
        typeset_offset: 0,
        constraint_offset: 0,
    },
    // IaddCin: fixed_results=1, use_typevar_operand=true, requires_typevar_operand=false, fixed_values=3
    // Constraints=['Same', 'Same', 'Same', 'Concrete(ir::types::I8)']
    // Polymorphic over TypeSet(lanes={1}, ints={8, 16, 32, 64, 128})
    OpcodeConstraints {
        flags: 0x69,
        typeset_offset: 0,
        constraint_offset: 39,
    },
    // IaddCout: fixed_results=2, use_typevar_operand=true, requires_typevar_operand=false, fixed_values=2
    // Constraints=['Same', 'Concrete(ir::types::I8)', 'Same', 'Same']
    // Polymorphic over TypeSet(lanes={1}, ints={8, 16, 32, 64, 128})
    OpcodeConstraints {
        flags: 0x4a,
        typeset_offset: 0,
        constraint_offset: 41,
    },
    // IaddCarry: fixed_results=2, use_typevar_operand=true, requires_typevar_operand=false, fixed_values=3
    // Constraints=['Same', 'Concrete(ir::types::I8)', 'Same', 'Same', 'Concrete(ir::types::I8)']
    // Polymorphic over TypeSet(lanes={1}, ints={8, 16, 32, 64, 128})
    OpcodeConstraints {
        flags: 0x6a,
        typeset_offset: 0,
        constraint_offset: 41,
    },
    // UaddOverflowTrap: fixed_results=1, use_typevar_operand=true, requires_typevar_operand=false, fixed_values=2
    // Constraints=['Same', 'Same', 'Same']
    // Polymorphic over TypeSet(lanes={1}, ints={32, 64})
    OpcodeConstraints {
        flags: 0x49,
        typeset_offset: 7,
        constraint_offset: 0,
    },
    // IsubBin: fixed_results=1, use_typevar_operand=true, requires_typevar_operand=false, fixed_values=3
    // Constraints=['Same', 'Same', 'Same', 'Concrete(ir::types::I8)']
    // Polymorphic over TypeSet(lanes={1}, ints={8, 16, 32, 64, 128})
    OpcodeConstraints {
        flags: 0x69,
        typeset_offset: 0,
        constraint_offset: 39,
    },
    // IsubBout: fixed_results=2, use_typevar_operand=true, requires_typevar_operand=false, fixed_values=2
    // Constraints=['Same', 'Concrete(ir::types::I8)', 'Same', 'Same']
    // Polymorphic over TypeSet(lanes={1}, ints={8, 16, 32, 64, 128})
    OpcodeConstraints {
        flags: 0x4a,
        typeset_offset: 0,
        constraint_offset: 41,
    },
    // IsubBorrow: fixed_results=2, use_typevar_operand=true, requires_typevar_operand=false, fixed_values=3
    // Constraints=['Same', 'Concrete(ir::types::I8)', 'Same', 'Same', 'Concrete(ir::types::I8)']
    // Polymorphic over TypeSet(lanes={1}, ints={8, 16, 32, 64, 128})
    OpcodeConstraints {
        flags: 0x6a,
        typeset_offset: 0,
        constraint_offset: 41,
    },
    // Band: fixed_results=1, use_typevar_operand=true, requires_typevar_operand=false, fixed_values=2
    // Constraints=['Same', 'Same', 'Same']
    // Polymorphic over TypeSet(lanes={1, 2, 4, 8, 16, 32, 64, 128, 256}, ints={8, 16, 32, 64, 128}, floats={32, 64})
    OpcodeConstraints {
        flags: 0x49,
        typeset_offset: 15,
        constraint_offset: 0,
    },
    // Bor: fixed_results=1, use_typevar_operand=true, requires_typevar_operand=false, fixed_values=2
    // Constraints=['Same', 'Same', 'Same']
    // Polymorphic over TypeSet(lanes={1, 2, 4, 8, 16, 32, 64, 128, 256}, ints={8, 16, 32, 64, 128}, floats={32, 64})
    OpcodeConstraints {
        flags: 0x49,
        typeset_offset: 15,
        constraint_offset: 0,
    },
    // Bxor: fixed_results=1, use_typevar_operand=true, requires_typevar_operand=false, fixed_values=2
    // Constraints=['Same', 'Same', 'Same']
    // Polymorphic over TypeSet(lanes={1, 2, 4, 8, 16, 32, 64, 128, 256}, ints={8, 16, 32, 64, 128}, floats={32, 64})
    OpcodeConstraints {
        flags: 0x49,
        typeset_offset: 15,
        constraint_offset: 0,
    },
    // Bnot: fixed_results=1, use_typevar_operand=true, requires_typevar_operand=false, fixed_values=1
    // Constraints=['Same', 'Same']
    // Polymorphic over TypeSet(lanes={1, 2, 4, 8, 16, 32, 64, 128, 256}, ints={8, 16, 32, 64, 128}, floats={32, 64})
    OpcodeConstraints {
        flags: 0x29,
        typeset_offset: 15,
        constraint_offset: 0,
    },
    // BandNot: fixed_results=1, use_typevar_operand=true, requires_typevar_operand=false, fixed_values=2
    // Constraints=['Same', 'Same', 'Same']
    // Polymorphic over TypeSet(lanes={1, 2, 4, 8, 16, 32, 64, 128, 256}, ints={8, 16, 32, 64, 128}, floats={32, 64})
    OpcodeConstraints {
        flags: 0x49,
        typeset_offset: 15,
        constraint_offset: 0,
    },
    // BorNot: fixed_results=1, use_typevar_operand=true, requires_typevar_operand=false, fixed_values=2
    // Constraints=['Same', 'Same', 'Same']
    // Polymorphic over TypeSet(lanes={1, 2, 4, 8, 16, 32, 64, 128, 256}, ints={8, 16, 32, 64, 128}, floats={32, 64})
    OpcodeConstraints {
        flags: 0x49,
        typeset_offset: 15,
        constraint_offset: 0,
    },
    // BxorNot: fixed_results=1, use_typevar_operand=true, requires_typevar_operand=false, fixed_values=2
    // Constraints=['Same', 'Same', 'Same']
    // Polymorphic over TypeSet(lanes={1, 2, 4, 8, 16, 32, 64, 128, 256}, ints={8, 16, 32, 64, 128}, floats={32, 64})
    OpcodeConstraints {
        flags: 0x49,
        typeset_offset: 15,
        constraint_offset: 0,
    },
    // BandImm: fixed_results=1, use_typevar_operand=true, requires_typevar_operand=false, fixed_values=1
    // Constraints=['Same', 'Same']
    // Polymorphic over TypeSet(lanes={1}, ints={8, 16, 32, 64, 128})
    OpcodeConstraints {
        flags: 0x29,
        typeset_offset: 0,
        constraint_offset: 0,
    },
    // BorImm: fixed_results=1, use_typevar_operand=true, requires_typevar_operand=false, fixed_values=1
    // Constraints=['Same', 'Same']
    // Polymorphic over TypeSet(lanes={1}, ints={8, 16, 32, 64, 128})
    OpcodeConstraints {
        flags: 0x29,
        typeset_offset: 0,
        constraint_offset: 0,
    },
    // BxorImm: fixed_results=1, use_typevar_operand=true, requires_typevar_operand=false, fixed_values=1
    // Constraints=['Same', 'Same']
    // Polymorphic over TypeSet(lanes={1}, ints={8, 16, 32, 64, 128})
    OpcodeConstraints {
        flags: 0x29,
        typeset_offset: 0,
        constraint_offset: 0,
    },
    // Rotl: fixed_results=1, use_typevar_operand=true, requires_typevar_operand=false, fixed_values=2
    // Constraints=['Same', 'Same', 'Free(0)']
    // Polymorphic over TypeSet(lanes={1, 2, 4, 8, 16, 32, 64, 128, 256}, ints={8, 16, 32, 64, 128})
    OpcodeConstraints {
        flags: 0x49,
        typeset_offset: 13,
        constraint_offset: 46,
    },
    // Rotr: fixed_results=1, use_typevar_operand=true, requires_typevar_operand=false, fixed_values=2
    // Constraints=['Same', 'Same', 'Free(0)']
    // Polymorphic over TypeSet(lanes={1, 2, 4, 8, 16, 32, 64, 128, 256}, ints={8, 16, 32, 64, 128})
    OpcodeConstraints {
        flags: 0x49,
        typeset_offset: 13,
        constraint_offset: 46,
    },
    // RotlImm: fixed_results=1, use_typevar_operand=true, requires_typevar_operand=false, fixed_values=1
    // Constraints=['Same', 'Same']
    // Polymorphic over TypeSet(lanes={1, 2, 4, 8, 16, 32, 64, 128, 256}, ints={8, 16, 32, 64, 128})
    OpcodeConstraints {
        flags: 0x29,
        typeset_offset: 13,
        constraint_offset: 0,
    },
    // RotrImm: fixed_results=1, use_typevar_operand=true, requires_typevar_operand=false, fixed_values=1
    // Constraints=['Same', 'Same']
    // Polymorphic over TypeSet(lanes={1, 2, 4, 8, 16, 32, 64, 128, 256}, ints={8, 16, 32, 64, 128})
    OpcodeConstraints {
        flags: 0x29,
        typeset_offset: 13,
        constraint_offset: 0,
    },
    // Ishl: fixed_results=1, use_typevar_operand=true, requires_typevar_operand=false, fixed_values=2
    // Constraints=['Same', 'Same', 'Free(0)']
    // Polymorphic over TypeSet(lanes={1, 2, 4, 8, 16, 32, 64, 128, 256}, ints={8, 16, 32, 64, 128})
    OpcodeConstraints {
        flags: 0x49,
        typeset_offset: 13,
        constraint_offset: 46,
    },
    // Ushr: fixed_results=1, use_typevar_operand=true, requires_typevar_operand=false, fixed_values=2
    // Constraints=['Same', 'Same', 'Free(0)']
    // Polymorphic over TypeSet(lanes={1, 2, 4, 8, 16, 32, 64, 128, 256}, ints={8, 16, 32, 64, 128})
    OpcodeConstraints {
        flags: 0x49,
        typeset_offset: 13,
        constraint_offset: 46,
    },
    // Sshr: fixed_results=1, use_typevar_operand=true, requires_typevar_operand=false, fixed_values=2
    // Constraints=['Same', 'Same', 'Free(0)']
    // Polymorphic over TypeSet(lanes={1, 2, 4, 8, 16, 32, 64, 128, 256}, ints={8, 16, 32, 64, 128})
    OpcodeConstraints {
        flags: 0x49,
        typeset_offset: 13,
        constraint_offset: 46,
    },
    // IshlImm: fixed_results=1, use_typevar_operand=true, requires_typevar_operand=false, fixed_values=1
    // Constraints=['Same', 'Same']
    // Polymorphic over TypeSet(lanes={1, 2, 4, 8, 16, 32, 64, 128, 256}, ints={8, 16, 32, 64, 128})
    OpcodeConstraints {
        flags: 0x29,
        typeset_offset: 13,
        constraint_offset: 0,
    },
    // UshrImm: fixed_results=1, use_typevar_operand=true, requires_typevar_operand=false, fixed_values=1
    // Constraints=['Same', 'Same']
    // Polymorphic over TypeSet(lanes={1, 2, 4, 8, 16, 32, 64, 128, 256}, ints={8, 16, 32, 64, 128})
    OpcodeConstraints {
        flags: 0x29,
        typeset_offset: 13,
        constraint_offset: 0,
    },
    // SshrImm: fixed_results=1, use_typevar_operand=true, requires_typevar_operand=false, fixed_values=1
    // Constraints=['Same', 'Same']
    // Polymorphic over TypeSet(lanes={1, 2, 4, 8, 16, 32, 64, 128, 256}, ints={8, 16, 32, 64, 128})
    OpcodeConstraints {
        flags: 0x29,
        typeset_offset: 13,
        constraint_offset: 0,
    },
    // Bitrev: fixed_results=1, use_typevar_operand=true, requires_typevar_operand=false, fixed_values=1
    // Constraints=['Same', 'Same']
    // Polymorphic over TypeSet(lanes={1}, ints={8, 16, 32, 64, 128})
    OpcodeConstraints {
        flags: 0x29,
        typeset_offset: 0,
        constraint_offset: 0,
    },
    // Clz: fixed_results=1, use_typevar_operand=true, requires_typevar_operand=false, fixed_values=1
    // Constraints=['Same', 'Same']
    // Polymorphic over TypeSet(lanes={1}, ints={8, 16, 32, 64, 128})
    OpcodeConstraints {
        flags: 0x29,
        typeset_offset: 0,
        constraint_offset: 0,
    },
    // Cls: fixed_results=1, use_typevar_operand=true, requires_typevar_operand=false, fixed_values=1
    // Constraints=['Same', 'Same']
    // Polymorphic over TypeSet(lanes={1}, ints={8, 16, 32, 64, 128})
    OpcodeConstraints {
        flags: 0x29,
        typeset_offset: 0,
        constraint_offset: 0,
    },
    // Ctz: fixed_results=1, use_typevar_operand=true, requires_typevar_operand=false, fixed_values=1
    // Constraints=['Same', 'Same']
    // Polymorphic over TypeSet(lanes={1}, ints={8, 16, 32, 64, 128})
    OpcodeConstraints {
        flags: 0x29,
        typeset_offset: 0,
        constraint_offset: 0,
    },
    // Bswap: fixed_results=1, use_typevar_operand=true, requires_typevar_operand=false, fixed_values=1
    // Constraints=['Same', 'Same']
    // Polymorphic over TypeSet(lanes={1}, ints={16, 32, 64, 128})
    OpcodeConstraints {
        flags: 0x29,
        typeset_offset: 16,
        constraint_offset: 0,
    },
    // Popcnt: fixed_results=1, use_typevar_operand=true, requires_typevar_operand=false, fixed_values=1
    // Constraints=['Same', 'Same']
    // Polymorphic over TypeSet(lanes={1, 2, 4, 8, 16, 32, 64, 128, 256}, ints={8, 16, 32, 64, 128})
    OpcodeConstraints {
        flags: 0x29,
        typeset_offset: 13,
        constraint_offset: 0,
    },
    // Fcmp: fixed_results=1, use_typevar_operand=true, requires_typevar_operand=true, fixed_values=2
    // Constraints=['AsBool', 'Same', 'Same']
    // Polymorphic over TypeSet(lanes={1, 2, 4, 8, 16, 32, 64, 128, 256}, floats={32, 64})
    OpcodeConstraints {
        flags: 0x59,
        typeset_offset: 17,
        constraint_offset: 33,
    },
    // Fadd: fixed_results=1, use_typevar_operand=true, requires_typevar_operand=false, fixed_values=2
    // Constraints=['Same', 'Same', 'Same']
    // Polymorphic over TypeSet(lanes={1, 2, 4, 8, 16, 32, 64, 128, 256}, floats={32, 64})
    OpcodeConstraints {
        flags: 0x49,
        typeset_offset: 17,
        constraint_offset: 0,
    },
    // Fsub: fixed_results=1, use_typevar_operand=true, requires_typevar_operand=false, fixed_values=2
    // Constraints=['Same', 'Same', 'Same']
    // Polymorphic over TypeSet(lanes={1, 2, 4, 8, 16, 32, 64, 128, 256}, floats={32, 64})
    OpcodeConstraints {
        flags: 0x49,
        typeset_offset: 17,
        constraint_offset: 0,
    },
    // Fmul: fixed_results=1, use_typevar_operand=true, requires_typevar_operand=false, fixed_values=2
    // Constraints=['Same', 'Same', 'Same']
    // Polymorphic over TypeSet(lanes={1, 2, 4, 8, 16, 32, 64, 128, 256}, floats={32, 64})
    OpcodeConstraints {
        flags: 0x49,
        typeset_offset: 17,
        constraint_offset: 0,
    },
    // Fdiv: fixed_results=1, use_typevar_operand=true, requires_typevar_operand=false, fixed_values=2
    // Constraints=['Same', 'Same', 'Same']
    // Polymorphic over TypeSet(lanes={1, 2, 4, 8, 16, 32, 64, 128, 256}, floats={32, 64})
    OpcodeConstraints {
        flags: 0x49,
        typeset_offset: 17,
        constraint_offset: 0,
    },
    // Sqrt: fixed_results=1, use_typevar_operand=true, requires_typevar_operand=false, fixed_values=1
    // Constraints=['Same', 'Same']
    // Polymorphic over TypeSet(lanes={1, 2, 4, 8, 16, 32, 64, 128, 256}, floats={32, 64})
    OpcodeConstraints {
        flags: 0x29,
        typeset_offset: 17,
        constraint_offset: 0,
    },
    // Fma: fixed_results=1, use_typevar_operand=true, requires_typevar_operand=false, fixed_values=3
    // Constraints=['Same', 'Same', 'Same', 'Same']
    // Polymorphic over TypeSet(lanes={1, 2, 4, 8, 16, 32, 64, 128, 256}, floats={32, 64})
    OpcodeConstraints {
        flags: 0x69,
        typeset_offset: 17,
        constraint_offset: 29,
    },
    // Fneg: fixed_results=1, use_typevar_operand=true, requires_typevar_operand=false, fixed_values=1
    // Constraints=['Same', 'Same']
    // Polymorphic over TypeSet(lanes={1, 2, 4, 8, 16, 32, 64, 128, 256}, floats={32, 64})
    OpcodeConstraints {
        flags: 0x29,
        typeset_offset: 17,
        constraint_offset: 0,
    },
    // Fabs: fixed_results=1, use_typevar_operand=true, requires_typevar_operand=false, fixed_values=1
    // Constraints=['Same', 'Same']
    // Polymorphic over TypeSet(lanes={1, 2, 4, 8, 16, 32, 64, 128, 256}, floats={32, 64})
    OpcodeConstraints {
        flags: 0x29,
        typeset_offset: 17,
        constraint_offset: 0,
    },
    // Fcopysign: fixed_results=1, use_typevar_operand=true, requires_typevar_operand=false, fixed_values=2
    // Constraints=['Same', 'Same', 'Same']
    // Polymorphic over TypeSet(lanes={1, 2, 4, 8, 16, 32, 64, 128, 256}, floats={32, 64})
    OpcodeConstraints {
        flags: 0x49,
        typeset_offset: 17,
        constraint_offset: 0,
    },
    // Fmin: fixed_results=1, use_typevar_operand=true, requires_typevar_operand=false, fixed_values=2
    // Constraints=['Same', 'Same', 'Same']
    // Polymorphic over TypeSet(lanes={1, 2, 4, 8, 16, 32, 64, 128, 256}, floats={32, 64})
    OpcodeConstraints {
        flags: 0x49,
        typeset_offset: 17,
        constraint_offset: 0,
    },
    // FminPseudo: fixed_results=1, use_typevar_operand=true, requires_typevar_operand=false, fixed_values=2
    // Constraints=['Same', 'Same', 'Same']
    // Polymorphic over TypeSet(lanes={1, 2, 4, 8, 16, 32, 64, 128, 256}, floats={32, 64})
    OpcodeConstraints {
        flags: 0x49,
        typeset_offset: 17,
        constraint_offset: 0,
    },
    // Fmax: fixed_results=1, use_typevar_operand=true, requires_typevar_operand=false, fixed_values=2
    // Constraints=['Same', 'Same', 'Same']
    // Polymorphic over TypeSet(lanes={1, 2, 4, 8, 16, 32, 64, 128, 256}, floats={32, 64})
    OpcodeConstraints {
        flags: 0x49,
        typeset_offset: 17,
        constraint_offset: 0,
    },
    // FmaxPseudo: fixed_results=1, use_typevar_operand=true, requires_typevar_operand=false, fixed_values=2
    // Constraints=['Same', 'Same', 'Same']
    // Polymorphic over TypeSet(lanes={1, 2, 4, 8, 16, 32, 64, 128, 256}, floats={32, 64})
    OpcodeConstraints {
        flags: 0x49,
        typeset_offset: 17,
        constraint_offset: 0,
    },
    // Ceil: fixed_results=1, use_typevar_operand=true, requires_typevar_operand=false, fixed_values=1
    // Constraints=['Same', 'Same']
    // Polymorphic over TypeSet(lanes={1, 2, 4, 8, 16, 32, 64, 128, 256}, floats={32, 64})
    OpcodeConstraints {
        flags: 0x29,
        typeset_offset: 17,
        constraint_offset: 0,
    },
    // Floor: fixed_results=1, use_typevar_operand=true, requires_typevar_operand=false, fixed_values=1
    // Constraints=['Same', 'Same']
    // Polymorphic over TypeSet(lanes={1, 2, 4, 8, 16, 32, 64, 128, 256}, floats={32, 64})
    OpcodeConstraints {
        flags: 0x29,
        typeset_offset: 17,
        constraint_offset: 0,
    },
    // Trunc: fixed_results=1, use_typevar_operand=true, requires_typevar_operand=false, fixed_values=1
    // Constraints=['Same', 'Same']
    // Polymorphic over TypeSet(lanes={1, 2, 4, 8, 16, 32, 64, 128, 256}, floats={32, 64})
    OpcodeConstraints {
        flags: 0x29,
        typeset_offset: 17,
        constraint_offset: 0,
    },
    // Nearest: fixed_results=1, use_typevar_operand=true, requires_typevar_operand=false, fixed_values=1
    // Constraints=['Same', 'Same']
    // Polymorphic over TypeSet(lanes={1, 2, 4, 8, 16, 32, 64, 128, 256}, floats={32, 64})
    OpcodeConstraints {
        flags: 0x29,
        typeset_offset: 17,
        constraint_offset: 0,
    },
    // IsNull: fixed_results=1, use_typevar_operand=true, requires_typevar_operand=true, fixed_values=1
    // Constraints=['Concrete(ir::types::I8)', 'Same']
    // Polymorphic over TypeSet(lanes={1}, refs={32, 64})
    OpcodeConstraints {
        flags: 0x39,
        typeset_offset: 11,
        constraint_offset: 36,
    },
    // IsInvalid: fixed_results=1, use_typevar_operand=true, requires_typevar_operand=true, fixed_values=1
    // Constraints=['Concrete(ir::types::I8)', 'Same']
    // Polymorphic over TypeSet(lanes={1}, refs={32, 64})
    OpcodeConstraints {
        flags: 0x39,
        typeset_offset: 11,
        constraint_offset: 36,
    },
    // Bitcast: fixed_results=1, use_typevar_operand=false, requires_typevar_operand=false, fixed_values=1
    // Constraints=['Same', 'Free(5)']
    // Polymorphic over TypeSet(lanes={1, 2, 4, 8, 16, 32, 64, 128, 256}, ints={8, 16, 32, 64, 128}, floats={32, 64}, refs={32, 64})
    OpcodeConstraints {
        flags: 0x21,
        typeset_offset: 5,
        constraint_offset: 49,
    },
    // ScalarToVector: fixed_results=1, use_typevar_operand=false, requires_typevar_operand=false, fixed_values=1
    // Constraints=['Same', 'LaneOf']
    // Polymorphic over TypeSet(lanes={2, 4, 8, 16, 32, 64, 128, 256}, ints={8, 16, 32, 64, 128}, floats={32, 64})
    OpcodeConstraints {
        flags: 0x21,
        typeset_offset: 10,
        constraint_offset: 4,
    },
    // Bmask: fixed_results=1, use_typevar_operand=false, requires_typevar_operand=false, fixed_values=1
    // Constraints=['Same', 'Free(3)']
    // Polymorphic over TypeSet(lanes={1, 2, 4, 8, 16, 32, 64, 128, 256}, ints={8, 16, 32, 64, 128})
    OpcodeConstraints {
        flags: 0x21,
        typeset_offset: 3,
        constraint_offset: 51,
    },
    // Ireduce: fixed_results=1, use_typevar_operand=false, requires_typevar_operand=false, fixed_values=1
    // Constraints=['Same', 'Free(0)']
    // Polymorphic over TypeSet(lanes={1}, ints={8, 16, 32, 64, 128})
    OpcodeConstraints {
        flags: 0x21,
        typeset_offset: 0,
        constraint_offset: 27,
    },
    // Snarrow: fixed_results=1, use_typevar_operand=true, requires_typevar_operand=true, fixed_values=2
    // Constraints=['SplitLanes', 'Same', 'Same']
    // Polymorphic over TypeSet(lanes={2, 4, 8}, ints={16, 32, 64})
    OpcodeConstraints {
        flags: 0x59,
        typeset_offset: 18,
        constraint_offset: 53,
    },
    // Unarrow: fixed_results=1, use_typevar_operand=true, requires_typevar_operand=true, fixed_values=2
    // Constraints=['SplitLanes', 'Same', 'Same']
    // Polymorphic over TypeSet(lanes={2, 4, 8}, ints={16, 32, 64})
    OpcodeConstraints {
        flags: 0x59,
        typeset_offset: 18,
        constraint_offset: 53,
    },
    // Uunarrow: fixed_results=1, use_typevar_operand=true, requires_typevar_operand=true, fixed_values=2
    // Constraints=['SplitLanes', 'Same', 'Same']
    // Polymorphic over TypeSet(lanes={2, 4, 8}, ints={16, 32, 64})
    OpcodeConstraints {
        flags: 0x59,
        typeset_offset: 18,
        constraint_offset: 53,
    },
    // SwidenLow: fixed_results=1, use_typevar_operand=true, requires_typevar_operand=true, fixed_values=1
    // Constraints=['MergeLanes', 'Same']
    // Polymorphic over TypeSet(lanes={2, 4, 8, 16}, ints={8, 16, 32})
    OpcodeConstraints {
        flags: 0x39,
        typeset_offset: 19,
        constraint_offset: 56,
    },
    // SwidenHigh: fixed_results=1, use_typevar_operand=true, requires_typevar_operand=true, fixed_values=1
    // Constraints=['MergeLanes', 'Same']
    // Polymorphic over TypeSet(lanes={2, 4, 8, 16}, ints={8, 16, 32})
    OpcodeConstraints {
        flags: 0x39,
        typeset_offset: 19,
        constraint_offset: 56,
    },
    // UwidenLow: fixed_results=1, use_typevar_operand=true, requires_typevar_operand=true, fixed_values=1
    // Constraints=['MergeLanes', 'Same']
    // Polymorphic over TypeSet(lanes={2, 4, 8, 16}, ints={8, 16, 32})
    OpcodeConstraints {
        flags: 0x39,
        typeset_offset: 19,
        constraint_offset: 56,
    },
    // UwidenHigh: fixed_results=1, use_typevar_operand=true, requires_typevar_operand=true, fixed_values=1
    // Constraints=['MergeLanes', 'Same']
    // Polymorphic over TypeSet(lanes={2, 4, 8, 16}, ints={8, 16, 32})
    OpcodeConstraints {
        flags: 0x39,
        typeset_offset: 19,
        constraint_offset: 56,
    },
    // IaddPairwise: fixed_results=1, use_typevar_operand=true, requires_typevar_operand=false, fixed_values=2
    // Constraints=['Same', 'Same', 'Same']
    // Polymorphic over TypeSet(lanes={2, 4, 8, 16}, ints={8, 16, 32})
    OpcodeConstraints {
        flags: 0x49,
        typeset_offset: 19,
        constraint_offset: 0,
    },
    // Uextend: fixed_results=1, use_typevar_operand=false, requires_typevar_operand=false, fixed_values=1
    // Constraints=['Same', 'Free(0)']
    // Polymorphic over TypeSet(lanes={1, 2, 4, 8, 16, 32, 64, 128, 256}, ints={8, 16, 32, 64, 128})
    OpcodeConstraints {
        flags: 0x21,
        typeset_offset: 3,
        constraint_offset: 27,
    },
    // Sextend: fixed_results=1, use_typevar_operand=false, requires_typevar_operand=false, fixed_values=1
    // Constraints=['Same', 'Free(0)']
    // Polymorphic over TypeSet(lanes={1, 2, 4, 8, 16, 32, 64, 128, 256}, ints={8, 16, 32, 64, 128})
    OpcodeConstraints {
        flags: 0x21,
        typeset_offset: 3,
        constraint_offset: 27,
    },
    // Fpromote: fixed_results=1, use_typevar_operand=false, requires_typevar_operand=false, fixed_values=1
    // Constraints=['Same', 'Free(17)']
    // Polymorphic over TypeSet(lanes={1, 2, 4, 8, 16, 32, 64, 128, 256}, floats={32, 64})
    OpcodeConstraints {
        flags: 0x21,
        typeset_offset: 20,
        constraint_offset: 57,
    },
    // Fdemote: fixed_results=1, use_typevar_operand=false, requires_typevar_operand=false, fixed_values=1
    // Constraints=['Same', 'Free(17)']
    // Polymorphic over TypeSet(lanes={1, 2, 4, 8, 16, 32, 64, 128, 256}, floats={32, 64})
    OpcodeConstraints {
        flags: 0x21,
        typeset_offset: 20,
        constraint_offset: 57,
    },
    // Fvdemote: fixed_results=1, use_typevar_operand=false, requires_typevar_operand=false, fixed_values=1
    // Constraints=['Concrete(ir::types::F32X4)', 'Concrete(ir::types::F64X2)']
    OpcodeConstraints {
        flags: 0x21,
        typeset_offset: 255,
        constraint_offset: 59,
    },
    // FvpromoteLow: fixed_results=1, use_typevar_operand=false, requires_typevar_operand=false, fixed_values=1
    // Constraints=['Concrete(ir::types::F64X2)', 'Concrete(ir::types::F32X4)']
    OpcodeConstraints {
        flags: 0x21,
        typeset_offset: 255,
        constraint_offset: 60,
    },
    // FcvtToUint: fixed_results=1, use_typevar_operand=false, requires_typevar_operand=false, fixed_values=1
    // Constraints=['Same', 'Free(21)']
    // Polymorphic over TypeSet(lanes={1, 2, 4, 8, 16, 32, 64, 128, 256}, ints={8, 16, 32, 64, 128})
    OpcodeConstraints {
        flags: 0x21,
        typeset_offset: 3,
        constraint_offset: 62,
    },
    // FcvtToSint: fixed_results=1, use_typevar_operand=false, requires_typevar_operand=false, fixed_values=1
    // Constraints=['Same', 'Free(21)']
    // Polymorphic over TypeSet(lanes={1, 2, 4, 8, 16, 32, 64, 128, 256}, ints={8, 16, 32, 64, 128})
    OpcodeConstraints {
        flags: 0x21,
        typeset_offset: 3,
        constraint_offset: 62,
    },
    // FcvtToUintSat: fixed_results=1, use_typevar_operand=false, requires_typevar_operand=false, fixed_values=1
    // Constraints=['Same', 'Free(17)']
    // Polymorphic over TypeSet(lanes={1, 2, 4, 8, 16, 32, 64, 128, 256}, ints={8, 16, 32, 64, 128})
    OpcodeConstraints {
        flags: 0x21,
        typeset_offset: 3,
        constraint_offset: 57,
    },
    // FcvtToSintSat: fixed_results=1, use_typevar_operand=false, requires_typevar_operand=false, fixed_values=1
    // Constraints=['Same', 'Free(17)']
    // Polymorphic over TypeSet(lanes={1, 2, 4, 8, 16, 32, 64, 128, 256}, ints={8, 16, 32, 64, 128})
    OpcodeConstraints {
        flags: 0x21,
        typeset_offset: 3,
        constraint_offset: 57,
    },
    // FcvtFromUint: fixed_results=1, use_typevar_operand=false, requires_typevar_operand=false, fixed_values=1
    // Constraints=['Same', 'Free(3)']
    // Polymorphic over TypeSet(lanes={1, 2, 4, 8, 16, 32, 64, 128, 256}, floats={32, 64})
    OpcodeConstraints {
        flags: 0x21,
        typeset_offset: 20,
        constraint_offset: 51,
    },
    // FcvtFromSint: fixed_results=1, use_typevar_operand=false, requires_typevar_operand=false, fixed_values=1
    // Constraints=['Same', 'Free(3)']
    // Polymorphic over TypeSet(lanes={1, 2, 4, 8, 16, 32, 64, 128, 256}, floats={32, 64})
    OpcodeConstraints {
        flags: 0x21,
        typeset_offset: 20,
        constraint_offset: 51,
    },
    // FcvtLowFromSint: fixed_results=1, use_typevar_operand=false, requires_typevar_operand=false, fixed_values=1
    // Constraints=['Same', 'Free(3)']
    // Polymorphic over TypeSet(lanes={1, 2, 4, 8, 16, 32, 64, 128, 256}, floats={32, 64})
    OpcodeConstraints {
        flags: 0x21,
        typeset_offset: 20,
        constraint_offset: 51,
    },
    // Isplit: fixed_results=2, use_typevar_operand=true, requires_typevar_operand=true, fixed_values=1
    // Constraints=['HalfWidth', 'HalfWidth', 'Same']
    // Polymorphic over TypeSet(lanes={1, 2, 4, 8, 16, 32, 64, 128, 256}, ints={16, 32, 64, 128})
    OpcodeConstraints {
        flags: 0x3a,
        typeset_offset: 22,
        constraint_offset: 64,
    },
    // Iconcat: fixed_results=1, use_typevar_operand=true, requires_typevar_operand=true, fixed_values=2
    // Constraints=['DoubleWidth', 'Same', 'Same']
    // Polymorphic over TypeSet(lanes={1, 2, 4, 8, 16, 32, 64, 128, 256}, ints={8, 16, 32, 64})
    OpcodeConstraints {
        flags: 0x59,
        typeset_offset: 9,
        constraint_offset: 67,
    },
    // AtomicRmw: fixed_results=1, use_typevar_operand=false, requires_typevar_operand=false, fixed_values=2
    // Constraints=['Same', 'Free(1)', 'Same']
    // Polymorphic over TypeSet(lanes={1}, ints={8, 16, 32, 64})
    OpcodeConstraints {
        flags: 0x41,
        typeset_offset: 23,
        constraint_offset: 69,
    },
    // AtomicCas: fixed_results=1, use_typevar_operand=true, requires_typevar_operand=false, fixed_values=3
    // Constraints=['Same', 'Free(1)', 'Same', 'Same']
    // Polymorphic over TypeSet(lanes={1}, ints={8, 16, 32, 64})
    OpcodeConstraints {
        flags: 0x69,
        typeset_offset: 23,
        constraint_offset: 69,
    },
    // AtomicLoad: fixed_results=1, use_typevar_operand=false, requires_typevar_operand=false, fixed_values=1
    // Constraints=['Same', 'Free(1)']
    // Polymorphic over TypeSet(lanes={1}, ints={8, 16, 32, 64})
    OpcodeConstraints {
        flags: 0x21,
        typeset_offset: 23,
        constraint_offset: 12,
    },
    // AtomicStore: fixed_results=0, use_typevar_operand=true, requires_typevar_operand=true, fixed_values=2
    // Constraints=['Same', 'Free(1)']
    // Polymorphic over TypeSet(lanes={1}, ints={8, 16, 32, 64})
    OpcodeConstraints {
        flags: 0x58,
        typeset_offset: 23,
        constraint_offset: 12,
    },
    // Fence: fixed_results=0, use_typevar_operand=false, requires_typevar_operand=false, fixed_values=0
    // Constraints=[]
    OpcodeConstraints {
        flags: 0x00,
        typeset_offset: 255,
        constraint_offset: 0,
    },
    // ExtractVector: fixed_results=1, use_typevar_operand=true, requires_typevar_operand=true, fixed_values=1
    // Constraints=['DynamicToVector', 'Same']
    // Polymorphic over TypeSet(lanes={1}, ints={8, 16, 32, 64, 128}, floats={32, 64})
    OpcodeConstraints {
        flags: 0x39,
        typeset_offset: 24,
        constraint_offset: 73,
    },
];

// Table of value type sets.
const TYPE_SETS: [ir::instructions::ValueTypeSet; 25] = [
    ir::instructions::ValueTypeSet {
        // TypeSet(lanes={1}, ints={8, 16, 32, 64, 128})
        lanes: BitSet::<u16>(1),
        dynamic_lanes: BitSet::<u16>(0),
        ints: BitSet::<u8>(248),
        floats: BitSet::<u8>(0),
        refs: BitSet::<u8>(0),
    },
    ir::instructions::ValueTypeSet {
        // TypeSet(lanes={1}, ints={32, 64}, refs={32, 64})
        lanes: BitSet::<u16>(1),
        dynamic_lanes: BitSet::<u16>(0),
        ints: BitSet::<u8>(96),
        floats: BitSet::<u8>(0),
        refs: BitSet::<u8>(96),
    },
    ir::instructions::ValueTypeSet {
        // TypeSet(lanes={2, 4, 8, 16, 32, 64, 128, 256}, ints={8, 16, 32, 64, 128}, floats={32, 64})
        lanes: BitSet::<u16>(510),
        dynamic_lanes: BitSet::<u16>(510),
        ints: BitSet::<u8>(248),
        floats: BitSet::<u8>(96),
        refs: BitSet::<u8>(0),
    },
    ir::instructions::ValueTypeSet {
        // TypeSet(lanes={1, 2, 4, 8, 16, 32, 64, 128, 256}, ints={8, 16, 32, 64, 128})
        lanes: BitSet::<u16>(511),
        dynamic_lanes: BitSet::<u16>(0),
        ints: BitSet::<u8>(248),
        floats: BitSet::<u8>(0),
        refs: BitSet::<u8>(0),
    },
    ir::instructions::ValueTypeSet {
        // TypeSet(lanes={2, 4, 8, 16, 32, 64, 128, 256}, ints={8, 16, 32, 64, 128})
        lanes: BitSet::<u16>(510),
        dynamic_lanes: BitSet::<u16>(0),
        ints: BitSet::<u8>(248),
        floats: BitSet::<u8>(0),
        refs: BitSet::<u8>(0),
    },
    ir::instructions::ValueTypeSet {
        // TypeSet(lanes={1, 2, 4, 8, 16, 32, 64, 128, 256}, ints={8, 16, 32, 64, 128}, floats={32, 64}, refs={32, 64})
        lanes: BitSet::<u16>(511),
        dynamic_lanes: BitSet::<u16>(510),
        ints: BitSet::<u8>(248),
        floats: BitSet::<u8>(96),
        refs: BitSet::<u8>(96),
    },
    ir::instructions::ValueTypeSet {
        // TypeSet(lanes={1}, ints={16, 32, 64})
        lanes: BitSet::<u16>(1),
        dynamic_lanes: BitSet::<u16>(0),
        ints: BitSet::<u8>(112),
        floats: BitSet::<u8>(0),
        refs: BitSet::<u8>(0),
    },
    ir::instructions::ValueTypeSet {
        // TypeSet(lanes={1}, ints={32, 64})
        lanes: BitSet::<u16>(1),
        dynamic_lanes: BitSet::<u16>(0),
        ints: BitSet::<u8>(96),
        floats: BitSet::<u8>(0),
        refs: BitSet::<u8>(0),
    },
    ir::instructions::ValueTypeSet {
        // TypeSet(lanes={1}, ints={64})
        lanes: BitSet::<u16>(1),
        dynamic_lanes: BitSet::<u16>(0),
        ints: BitSet::<u8>(64),
        floats: BitSet::<u8>(0),
        refs: BitSet::<u8>(0),
    },
    ir::instructions::ValueTypeSet {
        // TypeSet(lanes={1, 2, 4, 8, 16, 32, 64, 128, 256}, ints={8, 16, 32, 64})
        lanes: BitSet::<u16>(511),
        dynamic_lanes: BitSet::<u16>(510),
        ints: BitSet::<u8>(120),
        floats: BitSet::<u8>(0),
        refs: BitSet::<u8>(0),
    },
    ir::instructions::ValueTypeSet {
        // TypeSet(lanes={2, 4, 8, 16, 32, 64, 128, 256}, ints={8, 16, 32, 64, 128}, floats={32, 64})
        lanes: BitSet::<u16>(510),
        dynamic_lanes: BitSet::<u16>(0),
        ints: BitSet::<u8>(248),
        floats: BitSet::<u8>(96),
        refs: BitSet::<u8>(0),
    },
    ir::instructions::ValueTypeSet {
        // TypeSet(lanes={1}, refs={32, 64})
        lanes: BitSet::<u16>(1),
        dynamic_lanes: BitSet::<u16>(0),
        ints: BitSet::<u8>(0),
        floats: BitSet::<u8>(0),
        refs: BitSet::<u8>(96),
    },
    ir::instructions::ValueTypeSet {
        // TypeSet(lanes={1, 2, 4, 8, 16, 32, 64, 128, 256}, ints={8, 16, 32, 64, 128}, floats={32, 64}, refs={32, 64})
        lanes: BitSet::<u16>(511),
        dynamic_lanes: BitSet::<u16>(0),
        ints: BitSet::<u8>(248),
        floats: BitSet::<u8>(96),
        refs: BitSet::<u8>(96),
    },
    ir::instructions::ValueTypeSet {
        // TypeSet(lanes={1, 2, 4, 8, 16, 32, 64, 128, 256}, ints={8, 16, 32, 64, 128})
        lanes: BitSet::<u16>(511),
        dynamic_lanes: BitSet::<u16>(510),
        ints: BitSet::<u8>(248),
        floats: BitSet::<u8>(0),
        refs: BitSet::<u8>(0),
    },
    ir::instructions::ValueTypeSet {
        // TypeSet(lanes={4, 8}, ints={16, 32})
        lanes: BitSet::<u16>(12),
        dynamic_lanes: BitSet::<u16>(0),
        ints: BitSet::<u8>(48),
        floats: BitSet::<u8>(0),
        refs: BitSet::<u8>(0),
    },
    ir::instructions::ValueTypeSet {
        // TypeSet(lanes={1, 2, 4, 8, 16, 32, 64, 128, 256}, ints={8, 16, 32, 64, 128}, floats={32, 64})
        lanes: BitSet::<u16>(511),
        dynamic_lanes: BitSet::<u16>(0),
        ints: BitSet::<u8>(248),
        floats: BitSet::<u8>(96),
        refs: BitSet::<u8>(0),
    },
    ir::instructions::ValueTypeSet {
        // TypeSet(lanes={1}, ints={16, 32, 64, 128})
        lanes: BitSet::<u16>(1),
        dynamic_lanes: BitSet::<u16>(0),
        ints: BitSet::<u8>(240),
        floats: BitSet::<u8>(0),
        refs: BitSet::<u8>(0),
    },
    ir::instructions::ValueTypeSet {
        // TypeSet(lanes={1, 2, 4, 8, 16, 32, 64, 128, 256}, floats={32, 64})
        lanes: BitSet::<u16>(511),
        dynamic_lanes: BitSet::<u16>(510),
        ints: BitSet::<u8>(0),
        floats: BitSet::<u8>(96),
        refs: BitSet::<u8>(0),
    },
    ir::instructions::ValueTypeSet {
        // TypeSet(lanes={2, 4, 8}, ints={16, 32, 64})
        lanes: BitSet::<u16>(14),
        dynamic_lanes: BitSet::<u16>(14),
        ints: BitSet::<u8>(112),
        floats: BitSet::<u8>(0),
        refs: BitSet::<u8>(0),
    },
    ir::instructions::ValueTypeSet {
        // TypeSet(lanes={2, 4, 8, 16}, ints={8, 16, 32})
        lanes: BitSet::<u16>(30),
        dynamic_lanes: BitSet::<u16>(30),
        ints: BitSet::<u8>(56),
        floats: BitSet::<u8>(0),
        refs: BitSet::<u8>(0),
    },
    ir::instructions::ValueTypeSet {
        // TypeSet(lanes={1, 2, 4, 8, 16, 32, 64, 128, 256}, floats={32, 64})
        lanes: BitSet::<u16>(511),
        dynamic_lanes: BitSet::<u16>(0),
        ints: BitSet::<u8>(0),
        floats: BitSet::<u8>(96),
        refs: BitSet::<u8>(0),
    },
    ir::instructions::ValueTypeSet {
        // TypeSet(lanes={1}, floats={32, 64})
        lanes: BitSet::<u16>(1),
        dynamic_lanes: BitSet::<u16>(0),
        ints: BitSet::<u8>(0),
        floats: BitSet::<u8>(96),
        refs: BitSet::<u8>(0),
    },
    ir::instructions::ValueTypeSet {
        // TypeSet(lanes={1, 2, 4, 8, 16, 32, 64, 128, 256}, ints={16, 32, 64, 128})
        lanes: BitSet::<u16>(511),
        dynamic_lanes: BitSet::<u16>(0),
        ints: BitSet::<u8>(240),
        floats: BitSet::<u8>(0),
        refs: BitSet::<u8>(0),
    },
    ir::instructions::ValueTypeSet {
        // TypeSet(lanes={1}, ints={8, 16, 32, 64})
        lanes: BitSet::<u16>(1),
        dynamic_lanes: BitSet::<u16>(0),
        ints: BitSet::<u8>(120),
        floats: BitSet::<u8>(0),
        refs: BitSet::<u8>(0),
    },
    ir::instructions::ValueTypeSet {
        // TypeSet(lanes={1}, ints={8, 16, 32, 64, 128}, floats={32, 64})
        lanes: BitSet::<u16>(1),
        dynamic_lanes: BitSet::<u16>(510),
        ints: BitSet::<u8>(248),
        floats: BitSet::<u8>(96),
        refs: BitSet::<u8>(0),
    },
];

// Table of operand constraint sequences.
const OPERAND_CONSTRAINTS: [OperandConstraint; 75] = [
    OperandConstraint::Same,
    OperandConstraint::Same,
    OperandConstraint::Same,
    OperandConstraint::Concrete(ir::types::I32),
    OperandConstraint::Same,
    OperandConstraint::LaneOf,
    OperandConstraint::Concrete(ir::types::I8X16),
    OperandConstraint::Concrete(ir::types::I8X16),
    OperandConstraint::Concrete(ir::types::I8X16),
    OperandConstraint::Same,
    OperandConstraint::Same,
    OperandConstraint::LaneOf,
    OperandConstraint::Same,
    OperandConstraint::Free(1),
    OperandConstraint::Concrete(ir::types::I64),
    OperandConstraint::Same,
    OperandConstraint::Concrete(ir::types::I64),
    OperandConstraint::Free(1),
    OperandConstraint::Concrete(ir::types::I16X8),
    OperandConstraint::Same,
    OperandConstraint::Concrete(ir::types::I32X4),
    OperandConstraint::Same,
    OperandConstraint::Concrete(ir::types::I64X2),
    OperandConstraint::Same,
    OperandConstraint::Free(7),
    OperandConstraint::Concrete(ir::types::F32),
    OperandConstraint::Concrete(ir::types::F64),
    OperandConstraint::Same,
    OperandConstraint::Free(0),
    OperandConstraint::Same,
    OperandConstraint::Same,
    OperandConstraint::Same,
    OperandConstraint::Same,
    OperandConstraint::AsBool,
    OperandConstraint::Same,
    OperandConstraint::Same,
    OperandConstraint::Concrete(ir::types::I8),
    OperandConstraint::Same,
    OperandConstraint::Free(10),
    OperandConstraint::Same,
    OperandConstraint::Same,
    OperandConstraint::Same,
    OperandConstraint::Concrete(ir::types::I8),
    OperandConstraint::Same,
    OperandConstraint::Same,
    OperandConstraint::Concrete(ir::types::I8),
    OperandConstraint::Same,
    OperandConstraint::Same,
    OperandConstraint::Free(0),
    OperandConstraint::Same,
    OperandConstraint::Free(5),
    OperandConstraint::Same,
    OperandConstraint::Free(3),
    OperandConstraint::SplitLanes,
    OperandConstraint::Same,
    OperandConstraint::Same,
    OperandConstraint::MergeLanes,
    OperandConstraint::Same,
    OperandConstraint::Free(17),
    OperandConstraint::Concrete(ir::types::F32X4),
    OperandConstraint::Concrete(ir::types::F64X2),
    OperandConstraint::Concrete(ir::types::F32X4),
    OperandConstraint::Same,
    OperandConstraint::Free(21),
    OperandConstraint::HalfWidth,
    OperandConstraint::HalfWidth,
    OperandConstraint::Same,
    OperandConstraint::DoubleWidth,
    OperandConstraint::Same,
    OperandConstraint::Same,
    OperandConstraint::Free(1),
    OperandConstraint::Same,
    OperandConstraint::Same,
    OperandConstraint::DynamicToVector,
    OperandConstraint::Same,
];
