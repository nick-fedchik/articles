/-
  Mechanized model of the СТЕС admission kernel K (src/kernel) and a
  machine-checked proof of its end-to-end ALLOW soundness and exactness
  theorems (`soundness`, `allow_iff_P`, `allow_prohibitions_defeated_by_permit`).

  Lean 4 core only (no Mathlib). The model mirrors the Go definitions:
  * four-valued Belnap states as `(pos, neg)` bits (`B4`);
  * `truthAnd` conjunction of signed literals;
  * explicit acyclic signed Horn evidence rules (`EvRule`), their executable
    closure pass (`evalClosure`) and clean-support pass (`evalClean`, Theorem B.7),
    AND the declarative fixed-point certificate equations (`verifyCert`) checked
    across all atoms `a < S.atoms`;
  * `relevantAtom`: unreferenced noise facts on atoms outside guards, evidence
    rules, and norms do not poison `fresh` or `verified`;
  * `cleanWitness`: an ALLOW verdict requires at least one applicable,
    undefeated permit rule whose condition literals are all cleanly supported;
  * the status switch and decision switch of `evidence.(*NormProgram).evaluatePolicy`,
    the interval-guard convention of `compose.IntervalGuard`, freshness,
    provenance, and the kernel verdict downgrade.
-/

namespace GateKernel

structure B4 where
  pos : Bool
  neg : Bool
deriving DecidableEq, Repr

def B4.T : B4 := ⟨true, false⟩
def B4.F : B4 := ⟨false, true⟩
def B4.U : B4 := ⟨false, false⟩
def B4.Bo : B4 := ⟨true, true⟩

/-- Go `TruthAnd`: (a & b & 2) | ((a | b) & 1). -/
def truthAnd (a b : B4) : B4 := ⟨a.pos && b.pos, a.neg || b.neg⟩
def join (a b : B4) : B4 := ⟨a.pos || b.pos, a.neg || b.neg⟩

structure Lit where
  atom : Nat
  negv : Bool
deriving DecidableEq, Repr

def litVal (cl : Nat → B4) (l : Lit) : B4 :=
  if l.negv then ⟨(cl l.atom).neg, (cl l.atom).pos⟩ else cl l.atom

def litPos (cl : Nat → B4) (l : Lit) : Bool := (litVal cl l).pos

/-- One signed Horn evidence rule `premises → concl`. -/
structure EvRule where
  premises : List Lit
  concl : Lit
deriving DecidableEq, Repr

structure Norm where
  id : Nat
  permit : Bool
  conds : List Lit
  exc : Option Lit
deriving DecidableEq, Repr

/-- Condition value exactly as Go folds it: start TRUE, combine with TruthAnd. -/
def condVal (cl : Nat → B4) (r : Norm) : B4 :=
  r.conds.foldl (fun acc l => truthAnd acc (litVal cl l)) B4.T

def excVal (cl : Nat → B4) (r : Norm) : B4 :=
  match r.exc with
  | none => B4.F
  | some l => litVal cl l

inductive Status
  | inactive | excepted | unknownCond | conflictCond | unknownExc | conflictExc | applicable
deriving DecidableEq, Repr

/-- The switch of evaluatePolicy (all norms enabled in K). -/
def status (cl : Nat → B4) (r : Norm) : Status :=
  let c := condVal cl r
  let e := excVal cl r
  if c = B4.F then .inactive
  else if e = B4.T then .excepted
  else if c = B4.U then .unknownCond
  else if c = B4.Bo then .conflictCond
  else if e = B4.U then .unknownExc
  else if e = B4.Bo then .conflictExc
  else .applicable

inductive Decision
  | allow | deny | deferUnknown | deferConflict | deferNoBasis
deriving DecidableEq, Repr

def isUnknown : Status → Bool
  | .unknownCond | .unknownExc => true
  | _ => false

def isConflict : Status → Bool
  | .conflictCond | .conflictExc => true
  | _ => false

def dominated (cl : Nat → B4) (dom : Nat → Nat → Bool) (rules : List Norm) (r : Norm) : Bool :=
  rules.any (fun h => status cl h == .applicable && dom h.id r.id)

def survives (cl : Nat → B4) (dom : Nat → Nat → Bool) (rules : List Norm) (permit : Bool) : Bool :=
  rules.any (fun r => status cl r == .applicable && r.permit == permit && !dominated cl dom rules r)

/-- The decision switch of evaluatePolicy. -/
def decision (cl : Nat → B4) (dom : Nat → Nat → Bool) (rules : List Norm) : Decision :=
  if rules.any (fun r => isConflict (status cl r)) then .deferConflict
  else if rules.any (fun r => isUnknown (status cl r)) then .deferUnknown
  else match survives cl dom rules true, survives cl dom rules false with
    | true, true => .deferConflict
    | true, false => .allow
    | false, true => .deny
    | false, false => .deferNoBasis

structure Guard where
  atom : Nat
  qty : Nat
  unit : Nat
  lo : Int
  hi : Int
deriving DecidableEq, Repr

structure Meas where
  qty : Nat
  unit : Nat
  lo : Int
  hi : Int
  sources : List Nat
  time : Int
deriving DecidableEq, Repr

structure Fact where
  lit : Lit
  sources : List Nat
  time : Int
deriving DecidableEq, Repr

/-- compose.IntervalGuard: unit mismatch U, disjoint F, overlap BOTH, contained T.
    An inverted measurement interval `m.hi < m.lo` is malformed: Go `validateInput`
    rejects it with an error, the model maps it to U, so neither can yield ALLOW. -/
def guardState (g : Guard) (m : Meas) : B4 :=
  if m.unit ≠ g.unit ∨ m.hi < m.lo then B4.U
  else if m.hi < g.lo ∨ m.lo > g.hi then B4.F
  else if m.lo < g.lo ∨ m.hi > g.hi then B4.Bo
  else B4.T

structure Spec where
  atoms : Nat
  evid : List EvRule
  norms : List Norm
  dom : Nat → Nat → Bool
  guards : List Guard
  trusted : List Nat
  maxAge : Nat

/-- Topological acyclicity of the signed Horn evidence theory on `0 .. S.atoms-1`. -/
def Acyclic (S : Spec) : Prop :=
  ∀ r ∈ S.evid, r.concl.atom < S.atoms ∧ ∀ p ∈ r.premises, p.atom < r.concl.atom

/-- Well-formedness of evidence and norm literals within `0 .. S.atoms-1`. -/
def WellFormed (S : Spec) : Prop :=
  Acyclic S ∧
  ∀ r ∈ S.norms, (∀ l ∈ r.conds, l.atom < S.atoms) ∧
    match r.exc with | some e => e.atom < S.atoms | none => True

structure Input where
  now : Int
  facts : List Fact
  meas : List Meas

def factAtom (I : Input) (a : Nat) : B4 :=
  ⟨I.facts.any (fun f => f.lit.atom == a && !f.lit.negv),
   I.facts.any (fun f => f.lit.atom == a && f.lit.negv)⟩

def guardJoin (I : Input) (g : Guard) : B4 :=
  ⟨I.meas.any (fun m => m.qty == g.qty && (guardState g m).pos),
   I.meas.any (fun m => m.qty == g.qty && (guardState g m).neg)⟩

def preAtoms (S : Spec) (I : Input) (a : Nat) : B4 :=
  if a < S.atoms then
    match S.guards.find? (fun g => g.atom == a) with
    | some g => guardJoin I g
    | none => factAtom I a
  else B4.U

/-- Right-hand side of the signed closure fixed-point equation at literal `(a, negv)`. -/
def stepClosureBit (evid : List EvRule) (pre cl : Nat → B4) (a : Nat) (negv : Bool) : Bool :=
  (litVal pre ⟨a, negv⟩).pos ||
    evid.any (fun r => r.concl.atom == a && r.concl.negv == negv &&
      r.premises.all (fun p => litPos cl p))

def stepClosureAtom (evid : List EvRule) (pre cl : Nat → B4) (a : Nat) : B4 :=
  ⟨stepClosureBit evid pre cl a false, stepClosureBit evid pre cl a true⟩

/-- Right-hand side of the clean-support fixed-point equation (Theorem B.7) at `(a, negv)`. -/
def stepCleanBit (evid : List EvRule) (pre cl clean : Nat → B4) (a : Nat) (negv : Bool) : Bool :=
  (litVal cl ⟨a, negv⟩ == B4.T) &&
    ((litVal pre ⟨a, negv⟩).pos ||
      evid.any (fun r => r.concl.atom == a && r.concl.negv == negv &&
        r.premises.all (fun p => litPos clean p)))

def stepCleanAtom (evid : List EvRule) (pre cl clean : Nat → B4) (a : Nat) : B4 :=
  ⟨stepCleanBit evid pre cl clean a false, stepCleanBit evid pre cl clean a true⟩

/-- Forward topological pass computing the least signed closure over `0 .. n-1`. -/
def runClosure (evid : List EvRule) (pre : Nat → B4) : Nat → (Nat → B4)
  | 0 => fun _ => B4.U
  | k + 1 =>
    let cl := runClosure evid pre k
    fun a => if a == k then stepClosureAtom evid pre cl k else cl a

def closure (S : Spec) (I : Input) : Nat → B4 :=
  runClosure S.evid (preAtoms S I) S.atoms

/-- Forward topological pass computing the clean-support fixpoint over `0 .. n-1`. -/
def runClean (evid : List EvRule) (pre cl : Nat → B4) : Nat → (Nat → B4)
  | 0 => fun _ => B4.U
  | k + 1 =>
    let clean := runClean evid pre cl k
    fun a => if a == k then stepCleanAtom evid pre cl clean k else clean a

def clean (S : Spec) (I : Input) : Nat → B4 :=
  runClean S.evid (preAtoms S I) (closure S I) S.atoms

/-- An atom is relevant iff it appears in a guard, evidence rule, or norm of `S`. -/
def relevantAtom (S : Spec) (a : Nat) : Bool :=
  S.guards.any (fun g => g.atom == a) ||
  S.evid.any (fun r => r.concl.atom == a || r.premises.any (fun p => p.atom == a)) ||
  S.norms.any (fun r => r.conds.any (fun c => c.atom == a) ||
    match r.exc with | some e => e.atom == a | none => false)

def freshT (S : Spec) (I : Input) (t : Int) : Bool := decide (I.now - S.maxAge ≤ t ∧ t ≤ I.now)

def verifiedS (S : Spec) (srcs : List Nat) : Bool := !srcs.isEmpty && srcs.all (fun s => S.trusted.contains s)

def fresh (S : Spec) (I : Input) : Bool :=
  I.facts.all (fun f => !relevantAtom S f.lit.atom || freshT S I f.time) &&
  I.meas.all (fun m => freshT S I m.time)

def verified (S : Spec) (I : Input) : Bool :=
  I.facts.all (fun f => !relevantAtom S f.lit.atom || verifiedS S f.sources) &&
  I.meas.all (fun m => verifiedS S m.sources)

def contained (S : Spec) (I : Input) : Bool :=
  S.guards.all (fun g => I.meas.any (fun m => m.qty == g.qty) &&
    I.meas.all (fun m => m.qty != g.qty || guardState g m == B4.T))

def anyState (S : Spec) (I : Input) (s : B4) : Bool :=
  S.guards.any (fun g => I.meas.any (fun m => m.qty == g.qty && guardState g m == s))

/-- True iff some applicable, undefeated permit rule has all condition literals cleanly supported. -/
def cleanWitness (cl clean : Nat → B4) (dom : Nat → Nat → Bool) (rules : List Norm) : Bool :=
  rules.any (fun r => status cl r == .applicable && r.permit == true &&
    !dominated cl dom rules r && r.conds.all (fun l => litPos clean l))

def mapVerdict (S : Spec) (I : Input) (dec : Decision) (cw fr vr ct : Bool) : Decision :=
  match dec with
  | .allow =>
      if anyState S I B4.F then .deny
      else if anyState S I B4.Bo || !cw then .deferConflict
      else if fr && vr && ct then .allow
      else .deferUnknown
  | d => d

/-- Kernel verdict: the normative decision, downgraded when ALLOW violates a postcondition. -/
def verdict (S : Spec) (I : Input) : Decision :=
  let cl := closure S I
  let cln := clean S I
  mapVerdict S I (decision cl S.dom S.norms)
    (cleanWitness cl cln S.dom S.norms) (fresh S I) (verified S I) (contained S I)

structure Certificate where
  closure : Nat → B4
  clean : Nat → B4
  decision : Decision
  cleanWitness : Bool
  fresh : Bool
  verified : Bool
  contained : Bool
  verdict : Decision

def evaluate (S : Spec) (I : Input) : Certificate :=
  let cl := closure S I
  let cln := clean S I
  let dec := decision cl S.dom S.norms
  let cw := cleanWitness cl cln S.dom S.norms
  let fr := fresh S I
  let vr := verified S I
  let ct := contained S I
  ⟨cl, cln, dec, cw, fr, vr, ct, mapVerdict S I dec cw fr vr ct⟩

/-- Declarative certificate verifier: checks the closure and clean-support fixed-point
    equations across `0 .. S.atoms-1` and all postcondition equations. -/
def verifyCert (S : Spec) (I : Input) (c : Certificate) : Bool :=
  let pre := preAtoms S I
  (List.range S.atoms).all (fun a =>
    c.closure a == stepClosureAtom S.evid pre c.closure a &&
    c.clean a == stepCleanAtom S.evid pre c.closure c.clean a) &&
  c.decision == decision c.closure S.dom S.norms &&
  c.cleanWitness == cleanWitness c.closure c.clean S.dom S.norms &&
  c.fresh == fresh S I &&
  c.verified == verified S I &&
  c.contained == contained S I &&
  c.verdict == mapVerdict S I c.decision c.cleanWitness c.fresh c.verified c.contained

/-! ### The safety property P of an ALLOW decision -/

def P1 (S : Spec) (I : Input) : Prop :=
  ∃ r ∈ S.norms, r.permit = true ∧
    (∀ l ∈ r.conds, litVal (closure S I) l = B4.T ∧ litPos (clean S I) l = true) ∧
    excVal (closure S I) r = B4.F ∧
    ∀ h ∈ S.norms, status (closure S I) h = .applicable → S.dom h.id r.id = false

def P2 (S : Spec) (I : Input) : Prop :=
  ∀ r ∈ S.norms, r.permit = false → status (closure S I) r = .applicable →
    ∃ h ∈ S.norms, status (closure S I) h = .applicable ∧ S.dom h.id r.id = true

def P3 (S : Spec) (I : Input) : Prop :=
  ∀ r ∈ S.norms, condVal (closure S I) r = B4.F ∨ excVal (closure S I) r = B4.T ∨
    (condVal (closure S I) r = B4.T ∧ excVal (closure S I) r = B4.F)

def P4 (S : Spec) (I : Input) : Prop :=
  ∀ g ∈ S.guards, (∃ m ∈ I.meas, m.qty = g.qty) ∧
    ∀ m ∈ I.meas, m.qty = g.qty → m.unit = g.unit ∧ g.lo ≤ m.lo ∧ m.lo ≤ m.hi ∧ m.hi ≤ g.hi

def P5 (S : Spec) (I : Input) : Prop :=
  (∀ f ∈ I.facts, relevantAtom S f.lit.atom = true → I.now - S.maxAge ≤ f.time ∧ f.time ≤ I.now) ∧
  (∀ m ∈ I.meas, I.now - S.maxAge ≤ m.time ∧ m.time ≤ I.now)

def P6 (S : Spec) (I : Input) : Prop :=
  (∀ f ∈ I.facts, relevantAtom S f.lit.atom = true → f.sources ≠ [] ∧ ∀ s ∈ f.sources, s ∈ S.trusted) ∧
  (∀ m ∈ I.meas, m.sources ≠ [] ∧ ∀ s ∈ m.sources, s ∈ S.trusted)

def P (S : Spec) (I : Input) : Prop := P1 S I ∧ P2 S I ∧ P3 S I ∧ P4 S I ∧ P5 S I ∧ P6 S I

/-! ### Lemmas on Acyclic Fixed-Point Uniqueness and Existence -/

theorem all_congr {α} {l : List α} {p q : α → Bool} (h : ∀ x ∈ l, p x = q x) :
    l.all p = l.all q := by
  induction l with
  | nil => rfl
  | cons x xs ih =>
    simp only [List.all_cons]
    rw [h x List.mem_cons_self, ih (fun y hy => h y (List.mem_cons_of_mem x hy))]

theorem any_congr {α} {l : List α} {p q : α → Bool} (h : ∀ x ∈ l, p x = q x) :
    l.any p = l.any q := by
  induction l with
  | nil => rfl
  | cons x xs ih =>
    simp only [List.any_cons]
    rw [h x List.mem_cons_self, ih (fun y hy => h y (List.mem_cons_of_mem x hy))]

theorem stepClosureBit_congr (S : Spec) (hac : Acyclic S) (pre cl₁ cl₂ : Nat → B4) (a : Nat) (negv : Bool)
    (heq : ∀ b < a, cl₁ b = cl₂ b) :
    stepClosureBit S.evid pre cl₁ a negv = stepClosureBit S.evid pre cl₂ a negv := by
  unfold stepClosureBit
  congr 1
  apply any_congr
  intro r hr
  by_cases ha : r.concl.atom == a
  · have hca : r.concl.atom = a := beq_iff_eq.mp ha
    simp only [ha, Bool.true_and]
    congr 1
    apply all_congr
    intro p hp
    have hlt : p.atom < a := hca ▸ (hac r hr).2 p hp
    simp [litPos, litVal, heq p.atom hlt]
  · simp [ha]

theorem stepClosureAtom_congr (S : Spec) (hac : Acyclic S) (pre cl₁ cl₂ : Nat → B4) (a : Nat)
    (heq : ∀ b < a, cl₁ b = cl₂ b) :
    stepClosureAtom S.evid pre cl₁ a = stepClosureAtom S.evid pre cl₂ a := by
  unfold stepClosureAtom
  rw [stepClosureBit_congr S hac pre cl₁ cl₂ a false heq,
      stepClosureBit_congr S hac pre cl₁ cl₂ a true heq]

theorem stepCleanBit_congr (S : Spec) (hac : Acyclic S) (pre cl₁ cl₂ clean₁ clean₂ : Nat → B4)
    (a : Nat) (negv : Bool) (hcl : cl₁ a = cl₂ a) (heq : ∀ b < a, clean₁ b = clean₂ b) :
    stepCleanBit S.evid pre cl₁ clean₁ a negv = stepCleanBit S.evid pre cl₂ clean₂ a negv := by
  unfold stepCleanBit
  have hlit : litVal cl₁ ⟨a, negv⟩ = litVal cl₂ ⟨a, negv⟩ := by simp [litVal, hcl]
  rw [hlit]
  congr 2
  apply any_congr
  intro r hr
  by_cases ha : r.concl.atom == a
  · have hca : r.concl.atom = a := beq_iff_eq.mp ha
    simp only [ha, Bool.true_and]
    congr 1
    apply all_congr
    intro p hp
    have hlt : p.atom < a := hca ▸ (hac r hr).2 p hp
    simp [litPos, litVal, heq p.atom hlt]
  · simp [ha]

theorem stepCleanAtom_congr (S : Spec) (hac : Acyclic S) (pre cl₁ cl₂ clean₁ clean₂ : Nat → B4)
    (a : Nat) (hcl : cl₁ a = cl₂ a) (heq : ∀ b < a, clean₁ b = clean₂ b) :
    stepCleanAtom S.evid pre cl₁ clean₁ a = stepCleanAtom S.evid pre cl₂ clean₂ a := by
  unfold stepCleanAtom
  rw [stepCleanBit_congr S hac pre cl₁ cl₂ clean₁ clean₂ a false hcl heq,
      stepCleanBit_congr S hac pre cl₁ cl₂ clean₁ clean₂ a true hcl heq]

theorem runClosure_stable (evid : List EvRule) (pre : Nat → B4) (a k : Nat) (hlt : a < k) :
    runClosure evid pre k a = runClosure evid pre (a + 1) a := by
  induction k with
  | zero => omega
  | succ m ih =>
    rcases Nat.lt_succ_iff_lt_or_eq.mp hlt with hlt' | rfl
    · have hne : a ≠ m := Nat.ne_of_lt hlt'
      simp [runClosure, hne, ih hlt']
    · rfl

theorem runClosure_eq_step (S : Spec) (hac : Acyclic S) (pre : Nat → B4) (a : Nat) (ha : a < S.atoms) :
    runClosure S.evid pre S.atoms a = stepClosureAtom S.evid pre (runClosure S.evid pre S.atoms) a := by
  rw [runClosure_stable S.evid pre a S.atoms ha]
  simp only [runClosure, beq_self_eq_true, ite_true]
  apply stepClosureAtom_congr S hac pre (runClosure S.evid pre a) (runClosure S.evid pre S.atoms) a
  intro b hb
  rw [runClosure_stable S.evid pre b a hb,
      runClosure_stable S.evid pre b S.atoms (Nat.lt_trans hb ha)]

theorem runClean_stable (evid : List EvRule) (pre cl : Nat → B4) (a k : Nat) (hlt : a < k) :
    runClean evid pre cl k a = runClean evid pre cl (a + 1) a := by
  induction k with
  | zero => omega
  | succ m ih =>
    rcases Nat.lt_succ_iff_lt_or_eq.mp hlt with hlt' | rfl
    · have hne : a ≠ m := Nat.ne_of_lt hlt'
      simp [runClean, hne, ih hlt']
    · rfl

theorem runClean_eq_step (S : Spec) (hac : Acyclic S) (pre cl : Nat → B4) (a : Nat) (ha : a < S.atoms) :
    runClean S.evid pre cl S.atoms a = stepCleanAtom S.evid pre cl (runClean S.evid pre cl S.atoms) a := by
  rw [runClean_stable S.evid pre cl a S.atoms ha]
  simp only [runClean, beq_self_eq_true, ite_true]
  apply stepCleanAtom_congr S hac pre cl cl (runClean S.evid pre cl a) (runClean S.evid pre cl S.atoms) a rfl
  intro b hb
  rw [runClean_stable S.evid pre cl b a hb,
      runClean_stable S.evid pre cl b S.atoms (Nat.lt_trans hb ha)]

theorem fixpoint_closure_unique (S : Spec) (I : Input) (hac : Acyclic S) (cl : Nat → B4)
    (hfp : ∀ a < S.atoms, cl a = stepClosureAtom S.evid (preAtoms S I) cl a) :
    ∀ n ≤ S.atoms, ∀ a < n, cl a = closure S I a := by
  intro n
  induction n with
  | zero => intro _ a ha; omega
  | succ k ih =>
    intro hle a ha
    have hk : k < S.atoms := hle
    rcases Nat.lt_succ_iff_lt_or_eq.mp ha with hlt | rfl
    · exact ih (Nat.le_of_lt hle) a hlt
    · rw [hfp a hk, closure, runClosure_eq_step S hac (preAtoms S I) a hk]
      exact stepClosureAtom_congr S hac (preAtoms S I) cl (closure S I) a (ih (Nat.le_of_lt hle))

theorem fixpoint_clean_unique (S : Spec) (I : Input) (hac : Acyclic S) (cl cln : Nat → B4)
    (hcl : ∀ a < S.atoms, cl a = closure S I a)
    (hfp : ∀ a < S.atoms, cln a = stepCleanAtom S.evid (preAtoms S I) cl cln a) :
    ∀ n ≤ S.atoms, ∀ a < n, cln a = clean S I a := by
  intro n
  induction n with
  | zero => intro _ a ha; omega
  | succ k ih =>
    intro hle a ha
    have hk : k < S.atoms := hle
    rcases Nat.lt_succ_iff_lt_or_eq.mp ha with hlt | rfl
    · exact ih (Nat.le_of_lt hle) a hlt
    · rw [hfp a hk, clean, runClean_eq_step S hac (preAtoms S I) (closure S I) a hk]
      exact stepCleanAtom_congr S hac (preAtoms S I) cl (closure S I) cln (clean S I) a
        (hcl a hk) (ih (Nat.le_of_lt hle))

/-! ### Lemmas -/

theorem foldl_truthAnd_pos (cl : Nat → B4) (ls : List Lit) (acc : B4) :
    (ls.foldl (fun a l => truthAnd a (litVal cl l)) acc).pos = (acc.pos && ls.all (fun l => (litVal cl l).pos)) := by
  induction ls generalizing acc with
  | nil => simp
  | cons x xs ih =>
    simp only [List.foldl]
    rw [ih]
    simp [truthAnd, Bool.and_assoc]

theorem foldl_truthAnd_neg (cl : Nat → B4) (ls : List Lit) (acc : B4) :
    (ls.foldl (fun a l => truthAnd a (litVal cl l)) acc).neg = (acc.neg || ls.any (fun l => (litVal cl l).neg)) := by
  induction ls generalizing acc with
  | nil => simp
  | cons x xs ih =>
    simp only [List.foldl]
    rw [ih]
    simp [truthAnd, Bool.or_assoc]

theorem cond_true_lits (cl : Nat → B4) (r : Norm) (h : condVal cl r = B4.T) :
    ∀ l ∈ r.conds, litVal cl l = B4.T := by
  intro l hl
  have hp := congrArg B4.pos h
  have hn := congrArg B4.neg h
  simp only [condVal] at hp hn
  rw [foldl_truthAnd_pos] at hp
  rw [foldl_truthAnd_neg] at hn
  simp [B4.T, List.all_eq_true, List.any_eq_false] at hp hn
  have h1 := hp l hl
  have h2 := hn l hl
  cases hv : litVal cl l with
  | mk p n => simp_all [B4.T]

theorem B4.cases4 (x : B4) : x = B4.U ∨ x = B4.F ∨ x = B4.T ∨ x = B4.Bo := by
  cases x with
  | mk p n => cases p <;> cases n <;> simp [B4.U, B4.F, B4.T, B4.Bo]

theorem status_applicable (cl : Nat → B4) (r : Norm) :
    status cl r = .applicable ↔ condVal cl r = B4.T ∧ excVal cl r = B4.F := by
  unfold status
  rcases B4.cases4 (condVal cl r) with hc | hc | hc | hc <;>
  rcases B4.cases4 (excVal cl r) with he | he | he | he <;>
  simp [hc, he, B4.U, B4.F, B4.T, B4.Bo]

theorem status_determined (cl : Nat → B4) (r : Norm) (hu : isUnknown (status cl r) = false)
    (hc : isConflict (status cl r) = false) :
    condVal cl r = B4.F ∨ excVal cl r = B4.T ∨ (condVal cl r = B4.T ∧ excVal cl r = B4.F) := by
  unfold status at hu hc
  rcases B4.cases4 (condVal cl r) with h1 | h1 | h1 | h1 <;>
  rcases B4.cases4 (excVal cl r) with h2 | h2 | h2 | h2 <;>
  simp_all [isUnknown, isConflict, B4.U, B4.F, B4.T, B4.Bo]

theorem decision_allow (cl : Nat → B4) (dom : Nat → Nat → Bool) (rules : List Norm)
    (h : decision cl dom rules = .allow) :
    rules.any (fun r => isConflict (status cl r)) = false ∧
    rules.any (fun r => isUnknown (status cl r)) = false ∧
    survives cl dom rules true = true ∧ survives cl dom rules false = false := by
  unfold decision at h
  by_cases h1 : rules.any (fun r => isConflict (status cl r)) = true
  · simp [h1] at h
  · by_cases h2 : rules.any (fun r => isUnknown (status cl r)) = true
    · simp [h1, h2] at h
    · simp only [h1, h2, Bool.false_eq_true, ite_false] at h
      cases h3 : survives cl dom rules true <;> cases h4 : survives cl dom rules false <;>
        simp [h3, h4] at h
      exact ⟨by simpa using h1, by simpa using h2, rfl, rfl⟩

theorem foldl_truthAnd_congr (cl₁ cl₂ : Nat → B4) (ls : List Lit) (acc : B4)
    (h : ∀ l ∈ ls, litVal cl₁ l = litVal cl₂ l) :
    ls.foldl (fun a l => truthAnd a (litVal cl₁ l)) acc =
    ls.foldl (fun a l => truthAnd a (litVal cl₂ l)) acc := by
  induction ls generalizing acc with
  | nil => rfl
  | cons x xs ih =>
    simp only [List.foldl]
    rw [h x List.mem_cons_self, ih _ (fun y hy => h y (List.mem_cons_of_mem x hy))]

theorem condVal_congr (S : Spec) (cl₁ cl₂ : Nat → B4) (r : Norm) (hr : r ∈ S.norms)
    (hwf : WellFormed S) (heq : ∀ a < S.atoms, cl₁ a = cl₂ a) :
    condVal cl₁ r = condVal cl₂ r := by
  unfold condVal
  apply foldl_truthAnd_congr
  intro l hl
  have hlt := (hwf.2 r hr).1 l hl
  simp [litVal, heq l.atom hlt]

theorem excVal_congr (S : Spec) (cl₁ cl₂ : Nat → B4) (r : Norm) (hr : r ∈ S.norms)
    (hwf : WellFormed S) (heq : ∀ a < S.atoms, cl₁ a = cl₂ a) :
    excVal cl₁ r = excVal cl₂ r := by
  have he := (hwf.2 r hr).2
  unfold excVal
  cases h : r.exc with
  | none => rfl
  | some l =>
    rw [h] at he
    simp [litVal, heq l.atom he]

theorem status_congr (S : Spec) (cl₁ cl₂ : Nat → B4) (r : Norm) (hr : r ∈ S.norms)
    (hwf : WellFormed S) (heq : ∀ a < S.atoms, cl₁ a = cl₂ a) :
    status cl₁ r = status cl₂ r := by
  unfold status
  rw [condVal_congr S cl₁ cl₂ r hr hwf heq, excVal_congr S cl₁ cl₂ r hr hwf heq]

theorem dominated_congr (S : Spec) (cl₁ cl₂ : Nat → B4) (r : Norm)
    (hwf : WellFormed S) (heq : ∀ a < S.atoms, cl₁ a = cl₂ a) :
    dominated cl₁ S.dom S.norms r = dominated cl₂ S.dom S.norms r := by
  unfold dominated
  apply any_congr
  intro h hh
  rw [status_congr S cl₁ cl₂ h hh hwf heq]

theorem survives_congr (S : Spec) (cl₁ cl₂ : Nat → B4) (permit : Bool)
    (hwf : WellFormed S) (heq : ∀ a < S.atoms, cl₁ a = cl₂ a) :
    survives cl₁ S.dom S.norms permit = survives cl₂ S.dom S.norms permit := by
  unfold survives
  apply any_congr
  intro r hr
  rw [status_congr S cl₁ cl₂ r hr hwf heq, dominated_congr S cl₁ cl₂ r hwf heq]

theorem decision_congr (S : Spec) (cl₁ cl₂ : Nat → B4)
    (hwf : WellFormed S) (heq : ∀ a < S.atoms, cl₁ a = cl₂ a) :
    decision cl₁ S.dom S.norms = decision cl₂ S.dom S.norms := by
  unfold decision
  have h1 : S.norms.any (fun r => isConflict (status cl₁ r)) =
            S.norms.any (fun r => isConflict (status cl₂ r)) :=
    any_congr (fun r hr => by rw [status_congr S cl₁ cl₂ r hr hwf heq])
  have h2 : S.norms.any (fun r => isUnknown (status cl₁ r)) =
            S.norms.any (fun r => isUnknown (status cl₂ r)) :=
    any_congr (fun r hr => by rw [status_congr S cl₁ cl₂ r hr hwf heq])
  rw [h1, h2, survives_congr S cl₁ cl₂ true hwf heq, survives_congr S cl₁ cl₂ false hwf heq]

theorem cleanWitness_congr (S : Spec) (cl₁ cl₂ cln₁ cln₂ : Nat → B4)
    (hwf : WellFormed S) (hcl : ∀ a < S.atoms, cl₁ a = cl₂ a)
    (hcln : ∀ a < S.atoms, cln₁ a = cln₂ a) :
    cleanWitness cl₁ cln₁ S.dom S.norms = cleanWitness cl₂ cln₂ S.dom S.norms := by
  unfold cleanWitness
  apply any_congr
  intro r hr
  rw [status_congr S cl₁ cl₂ r hr hwf hcl, dominated_congr S cl₁ cl₂ r hwf hcl]
  congr 1
  apply all_congr
  intro l hl
  have hlt := (hwf.2 r hr).1 l hl
  simp [litPos, litVal, hcln l.atom hlt]

theorem verifyCert_verdict_eq (S : Spec) (I : Input) (c : Certificate)
    (hwf : WellFormed S) (hv : verifyCert S I c = true) :
    c.verdict = verdict S I := by
  unfold verifyCert at hv
  simp only [Bool.and_eq_true, List.all_eq_true, List.mem_range, beq_iff_eq] at hv
  obtain ⟨⟨⟨⟨⟨⟨hfp, hdec⟩, hcw⟩, hfr⟩, hvr⟩, hct⟩, hverd⟩ := hv
  have hcl : ∀ a < S.atoms, c.closure a = closure S I a :=
    fixpoint_closure_unique S I hwf.1 c.closure (fun a ha => (hfp a ha).1) S.atoms (Nat.le_refl _)
  have hcln : ∀ a < S.atoms, c.clean a = clean S I a :=
    fixpoint_clean_unique S I hwf.1 c.closure c.clean hcl (fun a ha => (hfp a ha).2) S.atoms (Nat.le_refl _)
  rw [hverd, hdec, hcw, hfr, hvr, hct,
      decision_congr S c.closure (closure S I) hwf hcl,
      cleanWitness_congr S c.closure (closure S I) c.clean (clean S I) hwf hcl hcln]
  rfl

theorem evaluate_verifies (S : Spec) (I : Input) (hac : Acyclic S) :
    verifyCert S I (evaluate S I) = true := by
  unfold verifyCert evaluate
  simp only [Bool.and_eq_true, List.all_eq_true, List.mem_range, beq_iff_eq]
  refine ⟨⟨⟨⟨⟨⟨?_, trivial⟩, trivial⟩, trivial⟩, trivial⟩, trivial⟩, trivial⟩
  intro a ha
  exact ⟨runClosure_eq_step S hac (preAtoms S I) a ha,
         runClean_eq_step S hac (preAtoms S I) (closure S I) a ha⟩

theorem verifyCert_unique (S : Spec) (I : Input) (c₁ c₂ : Certificate)
    (hwf : WellFormed S) (h₁ : verifyCert S I c₁ = true) (h₂ : verifyCert S I c₂ = true) :
    (∀ a < S.atoms, c₁.closure a = c₂.closure a ∧ c₁.clean a = c₂.clean a) ∧
    c₁.decision = c₂.decision ∧ c₁.cleanWitness = c₂.cleanWitness ∧
    c₁.fresh = c₂.fresh ∧ c₁.verified = c₂.verified ∧
    c₁.contained = c₂.contained ∧ c₁.verdict = c₂.verdict := by
  have hv₁ := h₁
  have hv₂ := h₂
  unfold verifyCert at hv₁ hv₂
  simp only [Bool.and_eq_true, List.all_eq_true, List.mem_range, beq_iff_eq] at hv₁ hv₂
  obtain ⟨⟨⟨⟨⟨⟨hfp₁, hdec₁⟩, hcw₁⟩, hfr₁⟩, hvr₁⟩, hct₁⟩, hverd₁⟩ := hv₁
  obtain ⟨⟨⟨⟨⟨⟨hfp₂, hdec₂⟩, hcw₂⟩, hfr₂⟩, hvr₂⟩, hct₂⟩, hverd₂⟩ := hv₂
  have hcl₁ : ∀ a < S.atoms, c₁.closure a = closure S I a :=
    fixpoint_closure_unique S I hwf.1 c₁.closure (fun a ha => (hfp₁ a ha).1) S.atoms (Nat.le_refl _)
  have hcl₂ : ∀ a < S.atoms, c₂.closure a = closure S I a :=
    fixpoint_closure_unique S I hwf.1 c₂.closure (fun a ha => (hfp₂ a ha).1) S.atoms (Nat.le_refl _)
  have hcln₁ : ∀ a < S.atoms, c₁.clean a = clean S I a :=
    fixpoint_clean_unique S I hwf.1 c₁.closure c₁.clean hcl₁ (fun a ha => (hfp₁ a ha).2) S.atoms (Nat.le_refl _)
  have hcln₂ : ∀ a < S.atoms, c₂.clean a = clean S I a :=
    fixpoint_clean_unique S I hwf.1 c₂.closure c₂.clean hcl₂ (fun a ha => (hfp₂ a ha).2) S.atoms (Nat.le_refl _)
  have hcl12 : ∀ a < S.atoms, c₁.closure a = c₂.closure a := fun a ha => (hcl₁ a ha).trans (hcl₂ a ha).symm
  have hcln12 : ∀ a < S.atoms, c₁.clean a = c₂.clean a := fun a ha => (hcln₁ a ha).trans (hcln₂ a ha).symm
  have hd : c₁.decision = c₂.decision := by
    rw [hdec₁, hdec₂, decision_congr S c₁.closure c₂.closure hwf hcl12]
  have hcw : c₁.cleanWitness = c₂.cleanWitness := by
    rw [hcw₁, hcw₂, cleanWitness_congr S c₁.closure c₂.closure c₁.clean c₂.clean hwf hcl12 hcln12]
  have hfr : c₁.fresh = c₂.fresh := hfr₁.trans hfr₂.symm
  have hvr : c₁.verified = c₂.verified := hvr₁.trans hvr₂.symm
  have hct : c₁.contained = c₂.contained := hct₁.trans hct₂.symm
  refine ⟨fun a ha => ⟨hcl12 a ha, hcln12 a ha⟩, hd, hcw, hfr, hvr, hct, ?_⟩
  rw [hverd₁, hverd₂, hd, hcw, hfr, hvr, hct]

theorem verdict_allow (S : Spec) (I : Input) (h : verdict S I = .allow) :
    decision (closure S I) S.dom S.norms = .allow ∧
    cleanWitness (closure S I) (clean S I) S.dom S.norms = true ∧
    fresh S I = true ∧ verified S I = true ∧ contained S I = true := by
  unfold verdict mapVerdict at h
  cases hd : decision (closure S I) S.dom S.norms <;> simp only [hd] at h <;> try contradiction
  refine ⟨rfl, ?_⟩
  by_cases a : anyState S I B4.F = true
  · simp [a] at h
  · by_cases b : (anyState S I B4.Bo || !cleanWitness (closure S I) (clean S I) S.dom S.norms) = true
    · simp [a, b] at h
    · by_cases c : (fresh S I && verified S I && contained S I) = true
      · simp only [Bool.or_eq_true, Bool.not_eq_true', not_or, Bool.not_eq_false] at b
        simp only [Bool.and_eq_true] at c
        exact ⟨b.2, c.1.1, c.1.2, c.2⟩
      · simp [a, b, c] at h

theorem guard_true_contained (g : Guard) (m : Meas) (h : guardState g m = B4.T) :
    m.unit = g.unit ∧ g.lo ≤ m.lo ∧ m.lo ≤ m.hi ∧ m.hi ≤ g.hi := by
  unfold guardState at h
  by_cases hu : m.unit ≠ g.unit ∨ m.hi < m.lo
  · simp [hu, B4.U, B4.T] at h
  · have hu' : m.unit = g.unit ∧ m.lo ≤ m.hi := by
      refine ⟨?_, ?_⟩
      · exact Classical.not_not.mp (fun hne => hu (Or.inl hne))
      · exact Int.not_lt.mp (fun hlt => hu (Or.inr hlt))
    by_cases hd : m.hi < g.lo ∨ m.lo > g.hi
    · simp [hu, hd, B4.F, B4.T] at h
    · by_cases ho : m.lo < g.lo ∨ m.hi > g.hi
      · simp [hu, hd, ho, B4.Bo, B4.T] at h
      · refine ⟨hu'.1, ?_, hu'.2, ?_⟩ <;> omega

/-! ### Main theorem -/

theorem allow_sound (S : Spec) (I : Input) (h : verdict S I = .allow) : P S I := by
  obtain ⟨hd, hcw, hf, hv, hc⟩ := verdict_allow S I h
  obtain ⟨hnc, hnu, _hperm, hproh⟩ := decision_allow _ _ _ hd
  refine ⟨?p1, ?p2, ?p3, ?p4, ?p5, ?p6⟩
  case p1 =>
    unfold cleanWitness at hcw
    rw [List.any_eq_true] at hcw
    obtain ⟨r, hr, hx⟩ := hcw
    simp only [Bool.and_eq_true, beq_iff_eq, Bool.not_eq_true', List.all_eq_true] at hx
    obtain ⟨⟨⟨happ, hpm⟩, hnd⟩, hcln⟩ := hx
    obtain ⟨hcond, hexc⟩ := (status_applicable _ _).mp happ
    have htr := cond_true_lits _ _ hcond
    refine ⟨r, hr, hpm, fun l hl => ⟨htr l hl, hcln l hl⟩, hexc, ?_⟩
    intro h' hh' happ'
    unfold dominated at hnd
    rw [List.any_eq_false] at hnd
    have := hnd h' hh'
    simp [happ'] at this
    exact this
  case p2 =>
    intro r hr hpf happ
    unfold survives at hproh
    rw [List.any_eq_false] at hproh
    have := hproh r hr
    simp only [happ, hpf, beq_self_eq_true, Bool.true_and, Bool.not_eq_true'] at this
    unfold dominated at this
    simp only [Bool.not_eq_false] at this
    rw [List.any_eq_true] at this
    obtain ⟨h', hh', hx⟩ := this
    simp only [Bool.and_eq_true, beq_iff_eq] at hx
    exact ⟨h', hh', hx.1, hx.2⟩
  case p3 =>
    intro r hr
    rw [List.any_eq_false] at hnc hnu
    exact status_determined _ _ (by simpa using hnu r hr) (by simpa using hnc r hr)
  case p4 =>
    intro g hg
    unfold contained at hc
    rw [List.all_eq_true] at hc
    have hx := hc g hg
    simp only [Bool.and_eq_true, List.any_eq_true, List.all_eq_true, beq_iff_eq,
      Bool.or_eq_true, bne_iff_ne, ne_eq] at hx
    obtain ⟨⟨m0, hm0, hq0⟩, hall⟩ := hx
    refine ⟨⟨m0, hm0, hq0⟩, ?_⟩
    intro m hm hq
    rcases hall m hm with hne | hT
    · exact absurd hq hne
    · exact guard_true_contained g m hT
  case p5 =>
    unfold fresh at hf
    simp only [Bool.and_eq_true, List.all_eq_true, Bool.or_eq_true, Bool.not_eq_true',
      freshT, decide_eq_true_eq] at hf
    refine ⟨fun f hf' hrel => ?_, hf.2⟩
    rcases hf.1 f hf' with hnr | hok
    · simp [hrel] at hnr
    · exact hok
  case p6 =>
    unfold verified at hv
    simp only [Bool.and_eq_true, List.all_eq_true, Bool.or_eq_true, Bool.not_eq_true',
      verifiedS, List.isEmpty_eq_false_iff, List.contains_iff_mem] at hv
    refine ⟨fun f hf' hrel => ?_, fun m hm => ⟨(hv.2 m hm).1, (hv.2 m hm).2⟩⟩
    rcases hv.1 f hf' with hnr | hok
    · simp [hrel] at hnr
    · exact ⟨hok.1, hok.2⟩

/-- End-to-end soundness: any certificate accepted by the fixed-point equation
    verifier `verifyCert` with verdict ALLOW satisfies P. -/
theorem soundness (S : Spec) (I : Input) (c : Certificate)
    (hwf : WellFormed S) (hv : verifyCert S I c = true) (ha : c.verdict = .allow) : P S I := by
  have hv_eq := verifyCert_verdict_eq S I c hwf hv
  rw [hv_eq] at ha
  exact allow_sound S I ha

/-! ### Strengthening of P2 under a strict priority order

  CompileNormative rejects cyclic priorities and uses the transitive closure, so
  `dom` is a strict order on rule identities. Then every applicable prohibition
  is defeated by an applicable PERMIT, not merely by some applicable rule. -/

def StrictOrd (dom : Nat → Nat → Bool) : Prop :=
  (∀ a, dom a a = false) ∧ (∀ a b c, dom a b = true → dom b c = true → dom a c = true)

set_option linter.unusedSimpArgs false in
theorem countP_lt {α} (p q : α → Bool) : ∀ (l : List α), (∀ x ∈ l, p x = true → q x = true) →
    ∀ y ∈ l, q y = true → p y = false → l.countP p < l.countP q
  | [], _, y, hy, _, _ => absurd hy List.not_mem_nil
  | a :: l, hm, y, hy, hq, hp => by
    have hm' : ∀ x ∈ l, p x = true → q x = true := fun x hx => hm x (List.mem_cons_of_mem a hx)
    have hle : l.countP p ≤ l.countP q := List.countP_mono_left hm'
    rcases List.mem_cons.mp hy with rfl | hy'
    · simp [hq, hp]
      omega
    · have ih := countP_lt p q l hm' y hy' hq hp
      simp only [List.countP_cons]
      by_cases hpa : p a = true
      · simp [hpa, hm a List.mem_cons_self hpa]
        omega
      · cases hqa : q a <;> simp [hpa, hqa] <;> omega

theorem exists_maximal (cl : Nat → B4) (dom : Nat → Nat → Bool) (rules : List Norm)
    (ho : StrictOrd dom) : ∀ n, ∀ r ∈ rules, status cl r = .applicable →
    rules.countP (fun h => status cl h == .applicable && dom h.id r.id) ≤ n →
    dominated cl dom rules r = false ∨
      ∃ g ∈ rules, status cl g = .applicable ∧ dominated cl dom rules g = false ∧ dom g.id r.id = true := by
  intro n
  induction n with
  | zero =>
    intro r _ _ hcount
    left
    unfold dominated
    rw [List.any_eq_false]
    intro h hh hx
    have : 0 < rules.countP (fun h => status cl h == .applicable && dom h.id r.id) :=
      List.countP_pos_iff.mpr ⟨h, hh, hx⟩
    omega
  | succ n ih =>
    intro r hr happ hcount
    cases hd : dominated cl dom rules r with
    | false => exact Or.inl rfl
    | true =>
      right
      unfold dominated at hd
      rw [List.any_eq_true] at hd
      obtain ⟨h, hh, hx⟩ := hd
      simp only [Bool.and_eq_true, beq_iff_eq] at hx
      obtain ⟨happh, hdom⟩ := hx
      have hlt := countP_lt (fun x => status cl x == .applicable && dom x.id h.id)
        (fun x => status cl x == .applicable && dom x.id r.id) rules
        (by
          intro x _ hx
          simp only [Bool.and_eq_true, beq_iff_eq] at hx ⊢
          exact ⟨hx.1, ho.2 _ _ _ hx.2 hdom⟩)
        h hh (by simp [happh, hdom]) (by simp [ho.1])
      rcases ih h hh happh (by omega) with hmax | ⟨g, hg, happg, hmaxg, hgh⟩
      · exact ⟨h, hh, happh, hmax, hdom⟩
      · exact ⟨g, hg, happg, hmaxg, ho.2 _ _ _ hgh hdom⟩

theorem allow_prohibitions_defeated_by_permit (S : Spec) (I : Input) (ho : StrictOrd S.dom)
    (h : verdict S I = .allow) :
    ∀ r ∈ S.norms, r.permit = false → status (closure S I) r = .applicable →
      ∃ g ∈ S.norms, g.permit = true ∧ status (closure S I) g = .applicable ∧ S.dom g.id r.id = true := by
  obtain ⟨hd, -, -, -, -⟩ := verdict_allow S I h
  obtain ⟨-, -, -, hproh⟩ := decision_allow _ _ _ hd
  unfold survives at hproh
  rw [List.any_eq_false] at hproh
  intro r hr hpf happ
  rcases exists_maximal (closure S I) S.dom S.norms ho _ r hr happ (Nat.le_refl _) with hnd | ⟨g, hg, happg, hmax, hgr⟩
  · have := hproh r hr
    simp [happ, hpf, hnd] at this
  · refine ⟨g, hg, ?_, happg, hgr⟩
    cases hgp : g.permit with
    | true => rfl
    | false =>
      have := hproh g hg
      simp [happg, hgp, hmax] at this

/-- Book form of P2: every applicable prohibition is dominated by an applicable PERMIT. -/
def P2permit (S : Spec) (I : Input) : Prop :=
  ∀ r ∈ S.norms, r.permit = false → status (closure S I) r = .applicable →
    ∃ g ∈ S.norms, g.permit = true ∧ status (closure S I) g = .applicable ∧ S.dom g.id r.id = true

/-- Certificate-level soundness with the book form of P2. The extra premise
    `StrictOrd S.dom` is what Go `CompileNormative` is designed to guarantee
    (cycle rejection and transitive closure); that link is not mechanized here. -/
theorem soundness_permit (S : Spec) (I : Input) (c : Certificate)
    (hwf : WellFormed S) (ho : StrictOrd S.dom)
    (hv : verifyCert S I c = true) (ha : c.verdict = .allow) : P S I ∧ P2permit S I := by
  have hv_eq := verifyCert_verdict_eq S I c hwf hv
  rw [hv_eq] at ha
  exact ⟨allow_sound S I ha, allow_prohibitions_defeated_by_permit S I ho ha⟩

/-! ### Relative completeness: P forces ALLOW -/

theorem lits_true_cond (cl : Nat → B4) (r : Norm) (h : ∀ l ∈ r.conds, litVal cl l = B4.T) :
    condVal cl r = B4.T := by
  have hp : (condVal cl r).pos = true := by
    simp only [condVal]
    rw [foldl_truthAnd_pos]
    simp only [B4.T, Bool.true_and, List.all_eq_true]
    intro l hl
    simp [h l hl, B4.T]
  have hn : (condVal cl r).neg = false := by
    simp only [condVal]
    rw [foldl_truthAnd_neg]
    simp only [B4.T, Bool.false_or, List.any_eq_false]
    intro l hl
    simp [h l hl, B4.T]
  cases hc : condVal cl r with
  | mk p n => simp_all [B4.T]

theorem determined_status (cl : Nat → B4) (r : Norm)
    (h : condVal cl r = B4.F ∨ excVal cl r = B4.T ∨ (condVal cl r = B4.T ∧ excVal cl r = B4.F)) :
    isUnknown (status cl r) = false ∧ isConflict (status cl r) = false := by
  unfold status
  rcases B4.cases4 (condVal cl r) with h1 | h1 | h1 | h1 <;>
  rcases B4.cases4 (excVal cl r) with h2 | h2 | h2 | h2 <;>
  simp_all [isUnknown, isConflict, B4.U, B4.F, B4.T, B4.Bo]

theorem contained_guard_true (g : Guard) (m : Meas)
    (h : m.unit = g.unit ∧ g.lo ≤ m.lo ∧ m.lo ≤ m.hi ∧ m.hi ≤ g.hi) : guardState g m = B4.T := by
  unfold guardState
  have h1 : ¬(m.hi < g.lo ∨ m.lo > g.hi) := by omega
  have h2 : ¬(m.lo < g.lo ∨ m.hi > g.hi) := by omega
  have h3 : ¬(m.hi < m.lo) := by omega
  simp [h.1, h1, h2, h3]

theorem allow_complete (S : Spec) (I : Input) (hP : P S I) :
    verdict S I = .allow := by
  obtain ⟨⟨r, hr, hpm, hlits, hexc, hnd⟩, hp2, hp3, hp4, hp5, hp6⟩ := hP
  have happ : status (closure S I) r = .applicable :=
    (status_applicable _ _).mpr ⟨lits_true_cond _ _ (fun l hl => (hlits l hl).1), hexc⟩
  have hnc : S.norms.any (fun x => isConflict (status (closure S I) x)) = false := by
    rw [List.any_eq_false]
    intro x hx
    simp [(determined_status _ _ (hp3 x hx)).2]
  have hnu : S.norms.any (fun x => isUnknown (status (closure S I) x)) = false := by
    rw [List.any_eq_false]
    intro x hx
    simp [(determined_status _ _ (hp3 x hx)).1]
  have hdom_false : dominated (closure S I) S.dom S.norms r = false := by
    unfold dominated
    rw [List.any_eq_false]
    intro h hh hx
    simp only [Bool.and_eq_true, beq_iff_eq] at hx
    simp [hnd h hh hx.1] at hx
  have hsurvT : survives (closure S I) S.dom S.norms true = true := by
    unfold survives
    rw [List.any_eq_true]
    refine ⟨r, hr, ?_⟩
    simp [happ, hpm, hdom_false]
  have hcw : cleanWitness (closure S I) (clean S I) S.dom S.norms = true := by
    unfold cleanWitness
    rw [List.any_eq_true]
    refine ⟨r, hr, ?_⟩
    have hall : r.conds.all (fun l => litPos (clean S I) l) = true :=
      List.all_eq_true.mpr (fun l hl => (hlits l hl).2)
    simp [happ, hpm, hdom_false, hall]
  have hsurvF : survives (closure S I) S.dom S.norms false = false := by
    unfold survives
    rw [List.any_eq_false]
    intro x hx hc
    simp only [Bool.and_eq_true, beq_iff_eq, Bool.not_eq_true'] at hc
    obtain ⟨⟨hxa, hxp⟩, hxd⟩ := hc
    obtain ⟨h, hh, happh, hdom⟩ := hp2 x hx hxp hxa
    unfold dominated at hxd
    rw [List.any_eq_false] at hxd
    have := hxd h hh
    simp [happh, hdom] at this
  have hdec : decision (closure S I) S.dom S.norms = .allow := by
    unfold decision
    simp [hnc, hnu, hsurvT, hsurvF]
  have hT : ∀ g ∈ S.guards, ∀ m ∈ I.meas, m.qty = g.qty → guardState g m = B4.T :=
    fun g hg m hm hq => contained_guard_true g m ((hp4 g hg).2 m hm hq)
  have hnot : ∀ s, s ≠ B4.T → anyState S I s = false := by
    intro s hs
    unfold anyState
    rw [List.any_eq_false]
    intro g hg hx
    rw [List.any_eq_true] at hx
    obtain ⟨m, hm, hx⟩ := hx
    simp only [Bool.and_eq_true, beq_iff_eq] at hx
    exact hs (hx.2 ▸ hT g hg m hm hx.1)
  have hF := hnot B4.F (by simp [B4.F, B4.T])
  have hB := hnot B4.Bo (by simp [B4.Bo, B4.T])
  have hfresh : fresh S I = true := by
    unfold fresh
    simp only [Bool.and_eq_true, List.all_eq_true, Bool.or_eq_true, Bool.not_eq_true',
      freshT, decide_eq_true_eq]
    refine ⟨fun f hf => ?_, hp5.2⟩
    cases hrel : relevantAtom S f.lit.atom
    · exact Or.inl rfl
    · exact Or.inr (hp5.1 f hf hrel)
  have hver : verified S I = true := by
    unfold verified
    simp only [Bool.and_eq_true, List.all_eq_true, Bool.or_eq_true, Bool.not_eq_true',
      verifiedS, List.isEmpty_eq_false_iff, List.contains_iff_mem]
    refine ⟨fun f hf => ?_, fun m hm => ⟨(hp6.2 m hm).1, (hp6.2 m hm).2⟩⟩
    cases hrel : relevantAtom S f.lit.atom
    · exact Or.inl rfl
    · exact Or.inr ⟨(hp6.1 f hf hrel).1, (hp6.1 f hf hrel).2⟩
  have hcont : contained S I = true := by
    unfold contained
    rw [List.all_eq_true]
    intro g hg
    obtain ⟨⟨m0, hm0, hq0⟩, -⟩ := hp4 g hg
    simp only [Bool.and_eq_true, List.any_eq_true, List.all_eq_true, beq_iff_eq,
      Bool.or_eq_true, bne_iff_ne, ne_eq]
    refine ⟨⟨m0, hm0, hq0⟩, fun m hm => ?_⟩
    by_cases hq : m.qty = g.qty
    · exact Or.inr (hT g hg m hm hq)
    · exact Or.inl hq
  unfold verdict mapVerdict
  simp [hdec, hF, hB, hcw, hfresh, hver, hcont]

/-- K is exact for P: ALLOW iff P. P4 itself demands ordered measurement
    intervals, so no separate precondition on `I.meas` is needed. -/
theorem allow_iff_P (S : Spec) (I : Input) :
    verdict S I = .allow ↔ P S I :=
  ⟨allow_sound S I, allow_complete S I⟩

end GateKernel
