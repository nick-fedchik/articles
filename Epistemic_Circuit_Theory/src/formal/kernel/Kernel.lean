/-
  Abstract model of the TEOS admission kernel K (src/teoskernel) and a
  mechanized proof of its end-to-end ALLOW soundness theorem.

  Lean 4 core only (no Mathlib). The model mirrors the Go definitions:
  four-valued Belnap states as (pos, neg) bits, TruthAnd conjunction of signed
  literals, the status switch of evidence.(*NormProgram).evaluatePolicy, the
  decision switch, the interval-guard convention of teoscompose.IntervalGuard,
  freshness, provenance and the kernel verdict mapping.

  Abstractions (stated, not hidden):
  * numbers are Int (Go uses finite float64; only <, <=, = are used);
  * the evidence closure is an arbitrary function `closeFn` supplied with the
    spec (Go: least fixpoint of acyclic signed Horn rules); the theorem holds
    for every such function, so it does not depend on the closure algorithm;
  * the priority relation `dom` is an arbitrary Boolean relation (Go: the
    transitive closure of an acyclic edge list);
  * Verify is modelled as recomputation and equality with the certificate,
    which is what the Go verifier does field by field.
-/

namespace TeosKernel

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

structure Norm where
  id : Nat
  permit : Bool
  conds : List Lit
  exc : Option Lit

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

structure Meas where
  qty : Nat
  unit : Nat
  lo : Int
  hi : Int
  sources : List Nat
  time : Int

structure Fact where
  lit : Lit
  sources : List Nat
  time : Int

/-- teoscompose.IntervalGuard: unit mismatch U, disjoint F, overlap BOTH, contained T. -/
def guardState (g : Guard) (m : Meas) : B4 :=
  if m.unit ≠ g.unit then B4.U
  else if m.hi < g.lo ∨ m.lo > g.hi then B4.F
  else if m.lo < g.lo ∨ m.hi > g.hi then B4.Bo
  else B4.T

structure Spec where
  norms : List Norm
  dom : Nat → Nat → Bool
  guards : List Guard
  trusted : List Nat
  maxAge : Nat
  /-- Evidence closure of a pre-closure assignment (Go: least fixpoint). -/
  closeFn : (Nat → B4) → (Nat → B4)

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
  match S.guards.find? (fun g => g.atom == a) with
  | some g => guardJoin I g
  | none => factAtom I a

def closure (S : Spec) (I : Input) : Nat → B4 := S.closeFn (preAtoms S I)

def freshT (S : Spec) (I : Input) (t : Int) : Bool := decide (I.now - S.maxAge ≤ t ∧ t ≤ I.now)

def verifiedS (S : Spec) (srcs : List Nat) : Bool := !srcs.isEmpty && srcs.all (fun s => S.trusted.contains s)

def fresh (S : Spec) (I : Input) : Bool :=
  I.facts.all (fun f => freshT S I f.time) && I.meas.all (fun m => freshT S I m.time)

def verified (S : Spec) (I : Input) : Bool :=
  I.facts.all (fun f => verifiedS S f.sources) && I.meas.all (fun m => verifiedS S m.sources)

def contained (S : Spec) (I : Input) : Bool :=
  S.guards.all (fun g => I.meas.any (fun m => m.qty == g.qty) &&
    I.meas.all (fun m => m.qty != g.qty || guardState g m == B4.T))

def anyState (S : Spec) (I : Input) (s : B4) : Bool :=
  S.guards.any (fun g => I.meas.any (fun m => m.qty == g.qty && guardState g m == s))

/-- Kernel verdict: the normative decision, downgraded when ALLOW violates a postcondition. -/
def verdict (S : Spec) (I : Input) : Decision :=
  match decision (closure S I) S.dom S.norms with
  | .allow =>
      if anyState S I B4.F then .deny
      else if anyState S I B4.Bo then .deferConflict
      else if fresh S I && verified S I && contained S I then .allow
      else .deferUnknown
  | d => d

structure Certificate where
  decision : Decision
  fresh : Bool
  verified : Bool
  contained : Bool
  verdict : Decision
deriving DecidableEq, Repr

def evaluate (S : Spec) (I : Input) : Certificate :=
  ⟨decision (closure S I) S.dom S.norms, fresh S I, verified S I, contained S I, verdict S I⟩

def verify (S : Spec) (I : Input) (c : Certificate) : Bool := c == evaluate S I

/-! ### The safety property P of an ALLOW decision -/

def P1 (S : Spec) (I : Input) : Prop :=
  ∃ r ∈ S.norms, r.permit = true ∧ (∀ l ∈ r.conds, litVal (closure S I) l = B4.T) ∧
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
    ∀ m ∈ I.meas, m.qty = g.qty → m.unit = g.unit ∧ g.lo ≤ m.lo ∧ m.hi ≤ g.hi

def P5 (S : Spec) (I : Input) : Prop :=
  (∀ f ∈ I.facts, I.now - S.maxAge ≤ f.time ∧ f.time ≤ I.now) ∧
  (∀ m ∈ I.meas, I.now - S.maxAge ≤ m.time ∧ m.time ≤ I.now)

def P6 (S : Spec) (I : Input) : Prop :=
  (∀ f ∈ I.facts, f.sources ≠ [] ∧ ∀ s ∈ f.sources, s ∈ S.trusted) ∧
  (∀ m ∈ I.meas, m.sources ≠ [] ∧ ∀ s ∈ m.sources, s ∈ S.trusted)

def P (S : Spec) (I : Input) : Prop := P1 S I ∧ P2 S I ∧ P3 S I ∧ P4 S I ∧ P5 S I ∧ P6 S I

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

theorem verdict_allow (S : Spec) (I : Input) (h : verdict S I = .allow) :
    decision (closure S I) S.dom S.norms = .allow ∧ fresh S I = true ∧ verified S I = true ∧
      contained S I = true := by
  unfold verdict at h
  cases hd : decision (closure S I) S.dom S.norms <;> simp only [hd] at h <;> try contradiction
  refine ⟨rfl, ?_⟩
  by_cases a : anyState S I B4.F = true
  · simp [a] at h
  · by_cases b : anyState S I B4.Bo = true
    · simp [a, b] at h
    · by_cases c : (fresh S I && verified S I && contained S I) = true
      · simp only [Bool.and_eq_true] at c
        exact ⟨c.1.1, c.1.2, c.2⟩
      · simp [a, b, c] at h

theorem guard_true_contained (g : Guard) (m : Meas) (h : guardState g m = B4.T) :
    m.unit = g.unit ∧ g.lo ≤ m.lo ∧ m.hi ≤ g.hi := by
  unfold guardState at h
  by_cases hu : m.unit ≠ g.unit
  · simp [hu, B4.U, B4.T] at h
  · by_cases hd : m.hi < g.lo ∨ m.lo > g.hi
    · simp [hu, hd, B4.F, B4.T] at h
    · by_cases ho : m.lo < g.lo ∨ m.hi > g.hi
      · simp [hu, hd, ho, B4.Bo, B4.T] at h
      · refine ⟨Classical.not_not.mp hu, ?_, ?_⟩ <;> omega

/-! ### Main theorem -/

theorem allow_sound (S : Spec) (I : Input) (h : verdict S I = .allow) : P S I := by
  obtain ⟨hd, hf, hv, hc⟩ := verdict_allow S I h
  obtain ⟨hnc, hnu, hperm, hproh⟩ := decision_allow _ _ _ hd
  refine ⟨?p1, ?p2, ?p3, ?p4, ?p5, ?p6⟩
  case p1 =>
    unfold survives at hperm
    rw [List.any_eq_true] at hperm
    obtain ⟨r, hr, hx⟩ := hperm
    simp only [Bool.and_eq_true, beq_iff_eq, Bool.not_eq_true'] at hx
    obtain ⟨⟨happ, hpm⟩, hnd⟩ := hx
    obtain ⟨hcond, hexc⟩ := (status_applicable _ _).mp happ
    refine ⟨r, hr, hpm, cond_true_lits _ _ hcond, hexc, ?_⟩
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
    simp only [Bool.and_eq_true, List.all_eq_true, freshT, decide_eq_true_eq] at hf
    exact ⟨hf.1, hf.2⟩
  case p6 =>
    unfold verified at hv
    simp only [Bool.and_eq_true, List.all_eq_true, verifiedS, Bool.not_eq_true',
      List.isEmpty_eq_false_iff, List.contains_iff_mem] at hv
    refine ⟨fun f hf' => ⟨(hv.1 f hf').1, (hv.1 f hf').2⟩, fun m hm => ⟨(hv.2 m hm).1, (hv.2 m hm).2⟩⟩

/-- End-to-end soundness: an accepted ALLOW certificate implies P. -/
theorem soundness (S : Spec) (I : Input) (c : Certificate)
    (hv : verify S I c = true) (ha : c.verdict = .allow) : P S I := by
  unfold verify at hv
  have hc : c = evaluate S I := by simpa using hv
  subst hc
  exact allow_sound S I ha

/-- Determinism: exactly one certificate is accepted for each (S, I). -/
theorem verify_unique (S : Spec) (I : Input) (c₁ c₂ : Certificate)
    (h₁ : verify S I c₁ = true) (h₂ : verify S I c₂ = true) : c₁ = c₂ := by
  unfold verify at h₁ h₂
  simp at h₁ h₂
  rw [h₁, h₂]

theorem evaluate_verifies (S : Spec) (I : Input) : verify S I (evaluate S I) = true := by
  simp [verify]

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
  obtain ⟨hd, -, -, -⟩ := verdict_allow S I h
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

/-! ### Relative completeness: P forces ALLOW on well-formed measurements -/

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

theorem contained_guard_true (g : Guard) (m : Meas) (hwf : m.lo ≤ m.hi)
    (h : m.unit = g.unit ∧ g.lo ≤ m.lo ∧ m.hi ≤ g.hi) : guardState g m = B4.T := by
  unfold guardState
  have h1 : ¬(m.hi < g.lo ∨ m.lo > g.hi) := by omega
  have h2 : ¬(m.lo < g.lo ∨ m.hi > g.hi) := by omega
  simp [h.1, h1, h2]

theorem allow_complete (S : Spec) (I : Input) (hwf : ∀ m ∈ I.meas, m.lo ≤ m.hi) (hP : P S I) :
    verdict S I = .allow := by
  obtain ⟨⟨r, hr, hpm, hlits, hexc, hnd⟩, hp2, hp3, hp4, hp5, hp6⟩ := hP
  have happ : status (closure S I) r = .applicable :=
    (status_applicable _ _).mpr ⟨lits_true_cond _ _ hlits, hexc⟩
  have hnc : S.norms.any (fun x => isConflict (status (closure S I) x)) = false := by
    rw [List.any_eq_false]
    intro x hx
    simp [(determined_status _ _ (hp3 x hx)).2]
  have hnu : S.norms.any (fun x => isUnknown (status (closure S I) x)) = false := by
    rw [List.any_eq_false]
    intro x hx
    simp [(determined_status _ _ (hp3 x hx)).1]
  have hsurvT : survives (closure S I) S.dom S.norms true = true := by
    unfold survives
    rw [List.any_eq_true]
    refine ⟨r, hr, ?_⟩
    have : dominated (closure S I) S.dom S.norms r = false := by
      unfold dominated
      rw [List.any_eq_false]
      intro h hh hx
      simp only [Bool.and_eq_true, beq_iff_eq] at hx
      simp [hnd h hh hx.1] at hx
    simp [happ, hpm, this]
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
    fun g hg m hm hq => contained_guard_true g m (hwf m hm) ((hp4 g hg).2 m hm hq)
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
    simp only [Bool.and_eq_true, List.all_eq_true, freshT, decide_eq_true_eq]
    exact hp5
  have hver : verified S I = true := by
    unfold verified
    simp only [Bool.and_eq_true, List.all_eq_true, verifiedS, Bool.not_eq_true',
      List.isEmpty_eq_false_iff, List.contains_iff_mem]
    exact ⟨fun f hf => hp6.1 f hf, fun m hm => hp6.2 m hm⟩
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
  unfold verdict
  simp [hdec, hF, hB, hfresh, hver, hcont]

/-- On well-formed measurements K is exact for P: ALLOW iff P. -/
theorem allow_iff_P (S : Spec) (I : Input) (hwf : ∀ m ∈ I.meas, m.lo ≤ m.hi) :
    verdict S I = .allow ↔ P S I :=
  ⟨allow_sound S I, allow_complete S I hwf⟩

end TeosKernel
