/**
 * Call System Types
 * Aligned with the 8 Blockers Framework - Doctrine/Principle Hierarchy
 */

// ============================================================================
// Structural Forces - Meta Layer
// ============================================================================

export type StructuralForce =
	| 'architectural_fragility'
	| 'organizational_power_ownership'
	| 'strategy_measurement_distortion';

// ============================================================================
// Core Blocker Types - The 8 Doctrines
// ============================================================================

export type BlockerType =
	| 'architecture_not_built_for_production'
	| 'data_dependencies_break_downstream'
	| 'unclear_ownership_turf_conflicts'
	| 'governance_safety_constraints'
	| 'vendor_driven_complexity'
	| 'strategy_execution_mismatch'
	| 'capacity_constraints'
	| 'misaligned_success_metrics';

export type BlockerSeverity = 'low' | 'medium' | 'high' | 'critical';

// ============================================================================
// Principle Types - Evidence Patterns Within Each Doctrine
// ============================================================================

export type PrincipleType =
	// Architecture principles
	| 'prototype_to_production_breakdown'
	| 'pattern_blindness'
	| 'sequence_confusion'
	| 'structural_debt'
	// Data principles
	| 'data_friction'
	| 'infrastructure_drift'
	| 'lineage_collapse'
	// Ownership principles
	| 'ownership_ambiguity'
	| 'authority_vacuum'
	| 'invisible_politics'
	| 'fragmented_decision_pathways'
	| 'capability_fracture'
	// Governance principles
	| 'governance_paralysis'
	| 'cultural_misfit'
	| 'risk_amplification'
	// Vendor principles
	| 'authority_vacuum_vendor'
	| 'tool_first_thinking'
	| 'complexity_tax'
	| 'lock_in_blindness'
	// Strategy principles
	| 'leadership_drift'
	| 'illusion_of_progress'
	| 'pattern_blindness_strategic'
	// Capacity principles
	| 'capacity_collapse'
	| 'invisible_queue'
	| 'capability_fracture_capacity'
	| 'bandwidth_blindness'
	// Metrics principles
	| 'misaligned_success_definitions'
	| 'metric_drift'
	| 'kpi_conflict'
	| 'premature_roi_pressure';

export interface PrincipleEvidence {
	principle: PrincipleType;
	parentDoctrine: BlockerType;
	evidence: string; // Direct quote from conversation
	timestamp: number;
	confidence: number;
}

// ============================================================================
// Blocker Interface
// ============================================================================

export interface Blocker {
	id: string;
	type: BlockerType; // Doctrine level
	severity: BlockerSeverity;
	description: string;
	principleEvidence: PrincipleEvidence[]; // Evidence at principle level
	solution?: string;
	identifiedAt: number;
	confidence?: number;
}

// ============================================================================
// Doctrine-Principle Mapping
// ============================================================================

export const DOCTRINE_PRINCIPLE_MAP: Record<BlockerType, PrincipleType[]> = {
	architecture_not_built_for_production: [
		'prototype_to_production_breakdown',
		'pattern_blindness',
		'sequence_confusion',
		'structural_debt'
	],
	data_dependencies_break_downstream: ['data_friction', 'infrastructure_drift', 'lineage_collapse'],
	unclear_ownership_turf_conflicts: [
		'ownership_ambiguity',
		'authority_vacuum',
		'invisible_politics',
		'fragmented_decision_pathways',
		'capability_fracture'
	],
	governance_safety_constraints: ['governance_paralysis', 'cultural_misfit', 'risk_amplification'],
	vendor_driven_complexity: [
		'authority_vacuum_vendor',
		'tool_first_thinking',
		'complexity_tax',
		'lock_in_blindness'
	],
	strategy_execution_mismatch: [
		'leadership_drift',
		'illusion_of_progress',
		'pattern_blindness_strategic',
		'sequence_confusion' // Note: shared with architecture
	],
	capacity_constraints: [
		'capacity_collapse',
		'invisible_queue',
		'capability_fracture_capacity',
		'bandwidth_blindness'
	],
	misaligned_success_metrics: [
		'misaligned_success_definitions',
		'metric_drift',
		'kpi_conflict',
		'premature_roi_pressure'
	]
};

// ============================================================================
// Structural Force Mapping
// ============================================================================

export const DOCTRINE_FORCE_MAP: Record<BlockerType, StructuralForce> = {
	architecture_not_built_for_production: 'architectural_fragility',
	data_dependencies_break_downstream: 'architectural_fragility',
	vendor_driven_complexity: 'architectural_fragility',
	unclear_ownership_turf_conflicts: 'organizational_power_ownership',
	governance_safety_constraints: 'organizational_power_ownership',
	capacity_constraints: 'organizational_power_ownership',
	strategy_execution_mismatch: 'strategy_measurement_distortion',
	misaligned_success_metrics: 'strategy_measurement_distortion'
};

// ============================================================================
// Blocker Metadata - Full Doctrine Definitions
// ============================================================================

export interface BlockerMetadata {
	type: BlockerType;
	label: string;
	icon: string;
	shortDescription: string;
	structuralForce: string;
	corePattern: string;
	principles: PrincipleType[];
	triggerQuestions: string[];
	rootExposureQuestions: string[];
	evidenceSignals: string[];
	patternName: string;
	bridge: string;
	commonMisdiagnosis: string[];
}

export const BLOCKER_METADATA: Record<BlockerType, BlockerMetadata> = {
	architecture_not_built_for_production: {
		type: 'architecture_not_built_for_production',
		label: 'Architecture Not Built for Production',
		icon: '⚙️',
		shortDescription: "Pilots work, production doesn't",
		structuralForce: 'Architectural Fragility',
		corePattern:
			'Prototype-grade systems break under real load. Scaling exposes gaps demos never reveal.',
		principles: [
			'prototype_to_production_breakdown',
			'pattern_blindness',
			'sequence_confusion',
			'structural_debt'
		],
		triggerQuestions: [
			'Tell me about your most promising pilot.',
			'What happens when you try to scale it?',
			'How many AI projects have made it to production?'
		],
		rootExposureQuestions: [
			'Who handles the transition from pilot to production?',
			'What breaks at scale?',
			'How long has the gap between demo and deployment existed?'
		],
		evidenceSignals: [
			'It worked in the demo',
			'We keep rebuilding',
			'Integration is the hard part',
			'The environment is different in production',
			'We had to roll it back'
		],
		patternName: 'You have an architecture-production gap.',
		bridge: 'When architecture aligns with production reality, rollout becomes predictable.',
		commonMisdiagnosis: [
			'We need better engineers',
			'We need more time',
			'The vendor tool is the problem'
		]
	},
	data_dependencies_break_downstream: {
		type: 'data_dependencies_break_downstream',
		label: 'Data Dependencies That Break Downstream',
		icon: '📊',
		shortDescription: 'Data is unstable, unreliable, or slow',
		structuralForce: 'Architectural Fragility',
		corePattern: 'A single undocumented dependency cascades into system-wide fragility.',
		principles: ['data_friction', 'infrastructure_drift', 'lineage_collapse'],
		triggerQuestions: [
			'How reliable is your data pipeline day-to-day?',
			'Is the org confident in the data feeding your models?',
			'How much time does the team spend on data prep vs. actual model work?'
		],
		rootExposureQuestions: [
			'Who owns data readiness?',
			'What happens when datasets shift?',
			'Can anyone explain how data flows from source to model?'
		],
		evidenceSignals: [
			'The data team is behind again',
			'It was fine until we added more users',
			"I'm not sure where that field comes from",
			'Different teams use different definitions',
			'The model started drifting'
		],
		patternName: 'This is a data dependency pattern.',
		bridge: 'Once data friction is removed, everything accelerates.',
		commonMisdiagnosis: [
			'We need better data scientists',
			'We need a new platform',
			"It's a data quality issue"
		]
	},
	unclear_ownership_turf_conflicts: {
		type: 'unclear_ownership_turf_conflicts',
		label: 'Unclear Ownership and Turf Conflicts',
		icon: '👤',
		shortDescription: 'No one actually owns AI strategy',
		structuralForce: 'Organizational Power and Ownership',
		corePattern: '"Who is accountable" matters more than "is the model accurate."',
		principles: [
			'ownership_ambiguity',
			'authority_vacuum',
			'invisible_politics',
			'fragmented_decision_pathways',
			'capability_fracture'
		],
		triggerQuestions: [
			'Who owns AI strategy today?',
			'Who is accountable if this fails?',
			'How does a decision move from idea to approval?'
		],
		rootExposureQuestions: [
			'If you took 30 days off, who makes decisions?',
			'Who sets priorities when things conflict?',
			'Who can say yes? Who can say no?',
			'What do people say off the record?'
		],
		evidenceSignals: [
			'I thought they were handling that',
			'The consultants are running point',
			"We're aligned in principle...",
			"It's stuck in the approval process",
			'They never asked us',
			'Everyone has a piece of it',
			"We're waiting on..."
		],
		patternName: "You're dealing with an ownership clarity issue.",
		bridge: 'When ownership becomes structural, movement becomes instant.',
		commonMisdiagnosis: [
			"It's a communication problem",
			'We need better collaboration',
			'People are resistant to change'
		]
	},
	governance_safety_constraints: {
		type: 'governance_safety_constraints',
		label: 'Governance and Safety Constraints',
		icon: '⚖️',
		shortDescription: 'Legal or compliance is slowing everything down',
		structuralForce: 'Organizational Power and Ownership',
		corePattern: 'Responsible governance crosses into operational paralysis.',
		principles: ['governance_paralysis', 'cultural_misfit', 'risk_amplification'],
		triggerQuestions: [
			'How does governance interact with your AI initiative?',
			'What slows decisions down most?',
			'How long does risk review take compared to development?'
		],
		rootExposureQuestions: [
			'What does governance need to feel comfortable?',
			'Who owns sign-off?',
			'Are different governance groups giving conflicting requirements?'
		],
		evidenceSignals: [
			"We're waiting on legal",
			"That approach won't work here",
			'We need to review this again',
			'Compliance keeps adding requirements',
			'No one wants to approve it',
			'The risk committee...'
		],
		patternName: "You're dealing with governance paralysis.",
		bridge: 'When governance becomes predictable, momentum returns.',
		commonMisdiagnosis: [
			'The organization is risk-averse',
			'Compliance is blocking us',
			'We need to educate legal'
		]
	},
	vendor_driven_complexity: {
		type: 'vendor_driven_complexity',
		label: 'Vendor-Driven Complexity',
		icon: '🔗',
		shortDescription: 'External dependencies constrain your options',
		structuralForce: 'Architectural Fragility',
		corePattern: 'Vendor architectures hard-wire fragility into your environment.',
		principles: [
			'authority_vacuum_vendor',
			'tool_first_thinking',
			'complexity_tax',
			'lock_in_blindness'
		],
		triggerQuestions: [
			'How dependent is your AI strategy on specific vendors?',
			'What happens if you need to switch platforms?',
			'Who made the decision on your current tooling?'
		],
		rootExposureQuestions: [
			'What would it cost to move off your current platform?',
			'Are you building capabilities or renting them?',
			'How much has vendor pricing changed since you started?'
		],
		evidenceSignals: [
			"That's what the vendor recommended",
			"We're a [Platform] shop",
			'The pricing changed when we scaled',
			"We'd have to rebuild everything",
			'We must do it this way because the platform requires it',
			'The vendor controls the roadmap'
		],
		patternName: "You're dealing with vendor-driven complexity.",
		bridge: 'When you own your architecture decisions, strategic flexibility returns.',
		commonMisdiagnosis: [
			'We need better vendor management',
			'We chose the wrong tool',
			'We need to renegotiate the contract'
		]
	},
	strategy_execution_mismatch: {
		type: 'strategy_execution_mismatch',
		label: 'Strategy-Execution Mismatch',
		icon: '🎯',
		shortDescription: 'Teams hit milestones but nothing reaches production',
		structuralForce: 'Strategy and Measurement Distortion',
		corePattern: 'Teams hit every milestone and still fail to reach production.',
		principles: [
			'leadership_drift',
			'sequence_confusion',
			'pattern_blindness_strategic',
			'illusion_of_progress'
		],
		triggerQuestions: [
			'How would you describe the gap between your AI vision and current reality?',
			'What happens to prototypes after the demo?',
			'How do business leaders and technical leaders describe success differently?'
		],
		rootExposureQuestions: [
			'If I asked engineering, product, and the COO what success looks like, would they say the same thing?',
			'How many initiatives have produced sustained business value?',
			"What's the real reason the last project stalled?"
		],
		evidenceSignals: [
			'The board is excited about this',
			"We're on track for Q3",
			'This is just a technology project',
			'We completed 14 initiatives this quarter',
			'Lots of activity, hard to point to results',
			'Strategy says one thing, reality is different'
		],
		patternName: 'You have a strategy-execution gap.',
		bridge: 'When strategy and execution align on structural truth, progress becomes real.',
		commonMisdiagnosis: [
			'Execution is too slow',
			'We need better project management',
			'Leadership needs to be more patient'
		]
	},
	capacity_constraints: {
		type: 'capacity_constraints',
		label: 'Capacity Constraints',
		icon: '⚡',
		shortDescription: 'People are overloaded',
		structuralForce: 'Organizational Power and Ownership',
		corePattern:
			'One overloaded senior engineer caps your AI trajectory without ever signaling "no."',
		principles: [
			'capacity_collapse',
			'invisible_queue',
			'capability_fracture_capacity',
			'bandwidth_blindness'
		],
		triggerQuestions: [
			'Who are the key contributors to AI work?',
			'Is this on top of their normal workload?',
			'What happens when your best people are unavailable?'
		],
		rootExposureQuestions: [
			'What percentage of their time is actually available for AI?',
			'What gets deprioritized when AI needs attention?',
			'If [key person] left tomorrow, what would stall?'
		],
		evidenceSignals: [
			"Everyone's stretched thin",
			"It's on the roadmap for next quarter",
			'Only [name] knows how that works',
			'This is in addition to their normal job',
			"We're trying to hire",
			'The same five people are on every project'
		],
		patternName: 'You have a capacity constraint.',
		bridge: 'Once capacity becomes structured, execution stabilizes.',
		commonMisdiagnosis: [
			'We need to hire',
			"The team isn't prioritizing correctly",
			'People need to work harder'
		]
	},
	misaligned_success_metrics: {
		type: 'misaligned_success_metrics',
		label: 'Misaligned Success Metrics',
		icon: '📈',
		shortDescription: 'Different teams want different outcomes',
		structuralForce: 'Strategy and Measurement Distortion',
		corePattern: 'Metrics designed to "prove progress" often prevent production entirely.',
		principles: [
			'misaligned_success_definitions',
			'metric_drift',
			'kpi_conflict',
			'premature_roi_pressure'
		],
		triggerQuestions: [
			'How does each team define success for AI?',
			'What does the board see as success?',
			'What metrics are you being measured on?'
		],
		rootExposureQuestions: [
			'If I asked engineering, legal, and the COO, would they say the same thing?',
			'Are you measuring activity or outcomes?',
			'When was the last time metrics changed behavior in a way that hurt the initiative?'
		],
		evidenceSignals: [
			'It depends who you ask',
			'We hit all our milestones',
			'The dashboard is green',
			"What's the payback period?",
			'Accuracy is high but nothing shipped',
			"We're measuring the wrong things"
		],
		patternName: 'You have a metrics alignment problem.',
		bridge: 'When success is unified, progress jumps.',
		commonMisdiagnosis: [
			'We need better reporting',
			"The team doesn't understand the business",
			'We need to track more metrics'
		]
	}
};

// ============================================================================
// Suggestion Types - AI-Generated Guidance for Agent
// ============================================================================

export type SuggestionType =
	| 'question'
	| 'pattern_reveal'
	| 'bridge'
	| 'qualification'
	| 'objection_handling';

export type SuggestionPriority = 'low' | 'medium' | 'high' | 'critical';

export interface Suggestion {
	id: string;
	suggestionType: SuggestionType;
	title: string;
	description: string;
	script?: string; // Exact language the agent can use verbatim
	priority: SuggestionPriority;
	relatedBlockerId?: string;
	forAgent: boolean; // Always true for suggestions
	shownToProspect: boolean; // Always false for suggestions
	dismissed?: boolean;
	used?: boolean;
	timestamp?: number;
}

// Question suggestion - what to ask next
export interface QuestionSuggestion extends Suggestion {
	suggestionType: 'question';
	blockerType: BlockerType;
	questionCategory: 'trigger' | 'root_exposure' | 'qualification';
}

// Pattern reveal suggestion - when to name the blocker
export interface PatternRevealSuggestion extends Suggestion {
	suggestionType: 'pattern_reveal';
	blockerType: BlockerType;
	evidenceCount: number; // Number of evidence pieces collected
	readyToReveal: boolean;
}

// Bridge suggestion - how to transition to offer
export interface BridgeSuggestion extends Suggestion {
	suggestionType: 'bridge';
	blockersCovered: string[]; // Which blockers have been identified
	qualificationComplete: boolean;
	recommendedOffer: 'clarity_retainer' | 'production_accelerator' | 'enterprise_ai_office' | 'nurture' | 'not_qualified';
}

// Qualification suggestion - verify buying authority
export interface QualificationSuggestion extends Suggestion {
	suggestionType: 'qualification';
	qualificationArea: 'budget_authority' | 'timeline_urgency' | 'decision_power';
	qualified?: boolean;
}

// Objection handling suggestion
export interface ObjectionHandlingSuggestion extends Suggestion {
	suggestionType: 'objection_handling';
	objectionType: 'price' | 'timing' | 'need_approval' | 'internal_solve' | 'too_busy';
}

// ============================================================================
// Call Phase Tracking
// ============================================================================

export type CallPhase = 'open' | 'diagnostic' | 'insight' | 'bridge' | 'close';

export interface PhaseProgress {
	currentPhase: CallPhase;
	phaseStartTime: number;
	timeInPhase: number; // Seconds
	totalCallDuration: number; // Seconds
	phasesCompleted: CallPhase[];
	nextPhaseRecommendation?: string;
}

export interface CallPhaseMetadata {
	phase: CallPhase;
	name: string;
	description: string;
	targetDuration: number; // Suggested duration in minutes
	objectives: string[];
	completionCriteria: string[];
}

export const CALL_PHASE_METADATA: Record<CallPhase, CallPhaseMetadata> = {
	open: {
		phase: 'open',
		name: 'Frame Setting',
		description: 'Establish safety, authority, and neutral tone',
		targetDuration: 3,
		objectives: [
			'Set expectations for the call',
			'Create psychological safety',
			'Get current state context'
		],
		completionCriteria: [
			'Opening script delivered',
			'First question asked',
			'Prospect has shared current state'
		]
	},
	diagnostic: {
		phase: 'diagnostic',
		name: 'Diagnostic Deep Dive',
		description: 'Identify which of the 8 blockers is active',
		targetDuration: 15,
		objectives: [
			'Ask trigger questions',
			'Gather evidence from responses',
			'Probe with root exposure questions',
			'Qualify buying authority'
		],
		completionCriteria: [
			'At least 1 blocker identified with evidence',
			'Qualification questions asked',
			'Pattern is clear'
		]
	},
	insight: {
		phase: 'insight',
		name: 'Insight Delivery',
		description: 'Name the pattern and create the "aha" moment',
		targetDuration: 5,
		objectives: [
			'Name the blocker pattern clearly',
			'Explain why it matters',
			'Show why common solutions fail',
			'Hint at the right approach'
		],
		completionCriteria: [
			'Pattern named and explained',
			'Prospect acknowledges pattern',
			'Value of solving is clear'
		]
	},
	bridge: {
		phase: 'bridge',
		name: 'Natural Bridge to Offer',
		description: 'Present the three-tier path and handle objections',
		targetDuration: 5,
		objectives: [
			'Ask if they want help solving',
			'Present the three-tier path',
			'Ask which tier they align with',
			'Handle objections'
		],
		completionCriteria: [
			'Three tiers explained',
			'Tier alignment asked',
			'Objections addressed',
			'Next steps clear'
		]
	},
	close: {
		phase: 'close',
		name: 'Close or Nurture',
		description: 'Get commitment or set up follow-up',
		targetDuration: 2,
		objectives: [
			'Confirm next steps',
			'Schedule follow-up if needed',
			'Send follow-up materials'
		],
		completionCriteria: [
			'Commitment received or nurture path set',
			'Follow-up scheduled',
			'Materials sent'
		]
	}
};

// ============================================================================
// Call Scripts - Verbatim Language for Agent
// ============================================================================

export const CALL_SCRIPTS = {
	opening: {
		frameSet: `Thanks for making time. Here's how this works.
For the next 20 minutes or so, my job is to help you see which structural pattern is actually blocking your AI initiative.
You'll walk away with clarity either way.
If I can help you fix it, I'll tell you what that looks like. If not, I'll tell you that too.
Does that work?`,
		firstQuestion: `So tell me—what's the current state of your AI work, and what made you want to have this conversation now?`
	},

	qualification: {
		budgetAuthority: `If we identified a clear structural fix today, do you have authority to act on it?`,
		timelineUrgency: `What happens if this isn't resolved by Q2?`,
		decisionPower: `Who else is involved in decisions at this level?`
	},

	patternReveal: (blockerType: BlockerType) => {
		const metadata = BLOCKER_METADATA[blockerType];
		return `Here's what I'm seeing.

${metadata.patternName}

[Explain in their context using their language]

The reason this matters is that most organizations try to fix this by [COMMON WRONG APPROACH], which never works because [WHY IT FAILS].

The 12% who reach production do something different. They [WHAT SUCCESSFUL ONES DO].

Does that match what you're experiencing?`;
	},

	bridge: {
		naturalTransition: `So the real question is: do you want to try solving this internally, or does it make sense for us to remove the blockers together?`,

		protocolIntro: `There are three paths from here.

One — Take the Pattern Map and execute internally.

Two — We remove the blockers together and reach production in the next 90 days. That's the $228,000 Production Accelerator.

Three — The $48,000 Clarity Retainer gives you the complete structural map and 30-day execution plan. You run it yourself.

The Enterprise AI Office starting at $760,000 is only for organizations that want me embedded as their fractional Head of AI Production across the full portfolio.

Which of those three feels most aligned right now?`
	},

	objectionHandling: {
		needToThink: `Of course. What specifically do you need to think through? Is it the investment, the timing, or something about the approach?`,

		needApproval: `That makes sense. What do you think their main question will be?

Here's how I'd frame it: "We've identified the structural blocker. For $48K you get the full structural map, the removal sequence, and the 30-day plan."`,

		solveInternally: `You absolutely could. The question is: how long will that take, and what's the cost of another quarter of drift?

Most teams I work with tried solving this internally for 6-9 months before engaging at the $48K clarity tier. The pattern doesn't become easier to see from inside.`,

		priceHigh: `I get it. Let me ask—what's the cost of another quarter with no measurable progress? What does that do to your credibility, your team's morale, your board's confidence?

The price isn't the investment. The investment is the time you're losing right now.`,

		tooBusy: `That's exactly why clarity matters. You're busy because you're stuck. Clarity is what gets you unstuck.

The Clarity Retainer requires 90 minutes for intake, one review session, and you have your execution plan within 30 days.`
	},

	close: {
		ready: `Great. I'll send the engagement link for the tier we discussed. Once the agreement is signed, we begin immediately. Your Pattern Map and removal sequencing will start within 48 hours.`,

		needTime: `No pressure. The clarity from today is yours either way. When you're ready, reply to my email and we'll begin the Clarity Retainer or the Production Accelerator.`
	}
};

// ============================================================================
// Qualification Questions
// ============================================================================

export interface QualificationQuestion {
	area: 'budget_authority' | 'timeline_urgency' | 'decision_power';
	question: string;
	positiveSignals: string[];
	negativeSignals: string[];
}

export const QUALIFICATION_QUESTIONS: QualificationQuestion[] = [
	{
		area: 'budget_authority',
		question: 'If we identified a clear structural fix today, do you have authority to act on it?',
		positiveSignals: [
			'Yes, I can approve this',
			'I own this budget',
			'I make these decisions',
			'Mentions specific budget range'
		],
		negativeSignals: [
			'Need to check with...',
			'Would need approval from...',
			'Not sure about budget',
			'Have to go through procurement'
		]
	},
	{
		area: 'timeline_urgency',
		question: "What happens if this isn't resolved by Q2?",
		positiveSignals: [
			'Board is asking questions',
			'Need to show progress',
			'Under pressure',
			'Specific deadline mentioned'
		],
		negativeSignals: [
			'No specific timeline',
			'Just exploring',
			'No urgency mentioned',
			'Long evaluation cycles'
		]
	},
	{
		area: 'decision_power',
		question: 'Who else is involved in decisions at this level?',
		positiveSignals: [
			'Just me',
			'I make the final call',
			'Small approval process',
			'Can move quickly'
		],
		negativeSignals: [
			'Committee decides',
			'Long list of stakeholders',
			'Complex approval process',
			'Multiple layers of sign-off'
		]
	}
];

// ============================================================================
// Transcript Types
// ============================================================================

export interface TranscriptSegment {
	id: string;
	speaker: 'agent' | 'prospect';
	text: string;
	timestamp: number;
	confidence?: number;
	blockerEvidence?: boolean; // Marked as evidence for a blocker
	relatedBlockerId?: string;
}

// ============================================================================
// Readiness Analysis (For Prospect View)
// ============================================================================

export interface ReadinessAnalysis {
	overallScore: number; // 0-100
	dimensions: {
		strategicAlignment: number;
		dataFoundation: number;
		technicalCapability: number;
		organizationalReadiness: number;
		governanceRisk: number;
		changeManagement: number;
	};
	trend: 'improving' | 'stable' | 'declining';
	primaryBlockers: BlockerType[];
	recommendations: string[];
}

// ============================================================================
// Session Types
// ============================================================================

export interface CallSession {
	id: string;
	leadId?: string;
	attendeeName?: string;
	attendeeCompany?: string;
	status: 'scheduled' | 'active' | 'completed' | 'cancelled';
	startTime?: string;
	endTime?: string;
	durationSeconds?: number;
	phaseProgress?: PhaseProgress;
	blockersIdentified?: Blocker[];
	suggestionsMade?: Suggestion[];
	qualificationStatus?: {
		budgetAuthority: boolean | null;
		timelineUrgency: boolean | null;
		decisionPower: boolean | null;
	};
	notes?: string;
	outcome?: 'clarity_retainer' | 'production_accelerator' | 'enterprise_ai_office' | 'nurture' | 'not_qualified' | 'no_show';
	reportGenerated?: boolean;
	reportUrl?: string;
}

// ============================================================================
// WebSocket Message Types
// ============================================================================

export type WebSocketMessageType =
	| 'transcript'
	| 'blocker'
	| 'suggestion'
	| 'phase_change'
	| 'analysis_update'
	| 'report_ready'
	| 'qualification_update';

export interface WebSocketMessage {
	type: WebSocketMessageType;
	data: unknown;
	timestamp?: number;
}

export interface TranscriptMessage extends WebSocketMessage {
	type: 'transcript';
	data: {
		speaker: 'agent' | 'prospect';
		text: string;
		timestamp: number;
		confidence?: number;
	};
}

export interface BlockerMessage extends WebSocketMessage {
	type: 'blocker';
	data: Blocker;
}

export interface SuggestionMessage extends WebSocketMessage {
	type: 'suggestion';
	data: Suggestion;
}

export interface PhaseChangeMessage extends WebSocketMessage {
	type: 'phase_change';
	data: {
		phase: CallPhase;
		timeElapsed: number;
		guidance: string;
		nextSteps: string[];
	};
}

export interface AnalysisUpdateMessage extends WebSocketMessage {
	type: 'analysis_update';
	data: ReadinessAnalysis;
}

export interface QualificationUpdateMessage extends WebSocketMessage {
	type: 'qualification_update';
	data: {
		budgetAuthority?: boolean;
		timelineUrgency?: boolean;
		decisionPower?: boolean;
		notes?: string;
	};
}

// ============================================================================
// Store-Compatible Types (for WebSocket messages)
// ============================================================================

export interface BlockerData {
	id: string;
	type: BlockerType;
	severity: BlockerSeverity;
	description: string;
	evidence?: string;
	solution?: string;
	identifiedAt: number;
}

export interface SuggestionData {
	id: string;
	type: string;
	title: string;
	description: string;
	priority: SuggestionPriority;
	forAgent: boolean;
	relatedBlockerId?: string;
}

export type ServerMessageType =
	| 'transcript'
	| 'blocker'
	| 'suggestion'
	| 'session_status'
	| 'connection_ack'
	| 'error';

export interface ServerMessage {
	type: ServerMessageType;
	// Transcript fields
	id?: string;
	speaker?: 'agent' | 'prospect';
	text?: string;
	timestamp?: number;
	confidence?: number;
	// Blocker fields
	blocker?: {
		id: string;
		type: BlockerType;
		severity: BlockerSeverity;
		description: string;
		evidence?: string;
		solution?: string;
		timestamp?: number;
	};
	// Suggestion fields
	suggestion?: {
		id: string;
		title: string;
		description: string;
		priority: SuggestionPriority;
		forAgent: boolean;
		relatedBlockerId?: string;
	};
	// Session status fields
	duration?: number;
	// Error fields
	error?: string;
}

// ============================================================================
// Utility Functions
// ============================================================================

export function getBlockerMetadata(type: BlockerType): BlockerMetadata {
	return BLOCKER_METADATA[type];
}

export function getPhaseMetadata(phase: CallPhase): CallPhaseMetadata {
	return CALL_PHASE_METADATA[phase];
}

export function calculateCallProgress(session: CallSession): number {
	if (!session.phaseProgress) return 0;
	const phases: CallPhase[] = ['open', 'diagnostic', 'insight', 'bridge', 'close'];
	const currentIndex = phases.indexOf(session.phaseProgress.currentPhase);
	return ((currentIndex + 1) / phases.length) * 100;
}

export function isQualified(session: CallSession): boolean {
	const qual = session.qualificationStatus;
	if (!qual) return false;
	return (
		qual.budgetAuthority === true &&
		qual.timelineUrgency === true &&
		qual.decisionPower === true
	);
}

export function getPrinciplesForDoctrine(doctrine: BlockerType): PrincipleType[] {
	return DOCTRINE_PRINCIPLE_MAP[doctrine];
}

export function getDoctrineForPrinciple(principle: PrincipleType): BlockerType | undefined {
	for (const [doctrine, principles] of Object.entries(DOCTRINE_PRINCIPLE_MAP)) {
		if (principles.includes(principle)) {
			return doctrine as BlockerType;
		}
	}
	return undefined;
}

export function getStructuralForce(doctrine: BlockerType): StructuralForce {
	return DOCTRINE_FORCE_MAP[doctrine];
}
