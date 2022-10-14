// Code generated by "stringer"; DO NOT EDIT.

package clusterversion

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[invalidVersionKey - -1]
	_ = x[V22_1-0]
	_ = x[Start22_2-1]
	_ = x[LocalTimestamps-2]
	_ = x[PebbleFormatSplitUserKeysMarkedCompacted-3]
	_ = x[EnsurePebbleFormatVersionRangeKeys-4]
	_ = x[EnablePebbleFormatVersionRangeKeys-5]
	_ = x[TrigramInvertedIndexes-6]
	_ = x[RemoveGrantPrivilege-7]
	_ = x[MVCCRangeTombstones-8]
	_ = x[UpgradeSequenceToBeReferencedByID-9]
	_ = x[SampledStmtDiagReqs-10]
	_ = x[AddSSTableTombstones-11]
	_ = x[SystemPrivilegesTable-12]
	_ = x[EnablePredicateProjectionChangefeed-13]
	_ = x[AlterSystemSQLInstancesAddLocality-14]
	_ = x[SystemExternalConnectionsTable-15]
	_ = x[AlterSystemStatementStatisticsAddIndexRecommendations-16]
	_ = x[RoleIDSequence-17]
	_ = x[AddSystemUserIDColumn-18]
	_ = x[SystemUsersIDColumnIsBackfilled-19]
	_ = x[SetSystemUsersUserIDColumnNotNull-20]
	_ = x[SQLSchemaTelemetryScheduledJobs-21]
	_ = x[SchemaChangeSupportsCreateFunction-22]
	_ = x[DeleteRequestReturnKey-23]
	_ = x[PebbleFormatPrePebblev1Marked-24]
	_ = x[RoleOptionsTableHasIDColumn-25]
	_ = x[RoleOptionsIDColumnIsBackfilled-26]
	_ = x[SetRoleOptionsUserIDColumnNotNull-27]
	_ = x[UseDelRangeInGCJob-28]
	_ = x[WaitedForDelRangeInGCJob-29]
	_ = x[RangefeedUseOneStreamPerNode-30]
	_ = x[NoNonMVCCAddSSTable-31]
	_ = x[GCHintInReplicaState-32]
	_ = x[UpdateInvalidColumnIDsInSequenceBackReferences-33]
	_ = x[TTLDistSQL-34]
	_ = x[PrioritizeSnapshots-35]
	_ = x[EnableLeaseUpgrade-36]
	_ = x[SupportAssumeRoleAuth-37]
	_ = x[FixUserfileRelatedDescriptorCorruption-38]
	_ = x[V22_2-39]
	_ = x[V23_1Start-40]
	_ = x[TenantNames-41]
}

const _Key_name = "invalidVersionKeyV22_1Start22_2LocalTimestampsPebbleFormatSplitUserKeysMarkedCompactedEnsurePebbleFormatVersionRangeKeysEnablePebbleFormatVersionRangeKeysTrigramInvertedIndexesRemoveGrantPrivilegeMVCCRangeTombstonesUpgradeSequenceToBeReferencedByIDSampledStmtDiagReqsAddSSTableTombstonesSystemPrivilegesTableEnablePredicateProjectionChangefeedAlterSystemSQLInstancesAddLocalitySystemExternalConnectionsTableAlterSystemStatementStatisticsAddIndexRecommendationsRoleIDSequenceAddSystemUserIDColumnSystemUsersIDColumnIsBackfilledSetSystemUsersUserIDColumnNotNullSQLSchemaTelemetryScheduledJobsSchemaChangeSupportsCreateFunctionDeleteRequestReturnKeyPebbleFormatPrePebblev1MarkedRoleOptionsTableHasIDColumnRoleOptionsIDColumnIsBackfilledSetRoleOptionsUserIDColumnNotNullUseDelRangeInGCJobWaitedForDelRangeInGCJobRangefeedUseOneStreamPerNodeNoNonMVCCAddSSTableGCHintInReplicaStateUpdateInvalidColumnIDsInSequenceBackReferencesTTLDistSQLPrioritizeSnapshotsEnableLeaseUpgradeSupportAssumeRoleAuthFixUserfileRelatedDescriptorCorruptionV22_2V23_1StartTenantNames"

var _Key_index = [...]uint16{0, 17, 22, 31, 46, 86, 120, 154, 176, 196, 215, 248, 267, 287, 308, 343, 377, 407, 460, 474, 495, 526, 559, 590, 624, 646, 675, 702, 733, 766, 784, 808, 836, 855, 875, 921, 931, 950, 968, 989, 1027, 1032, 1042, 1053}

func (i Key) String() string {
	i -= -1
	if i < 0 || i >= Key(len(_Key_index)-1) {
		return "Key(" + strconv.FormatInt(int64(i+-1), 10) + ")"
	}
	return _Key_name[_Key_index[i]:_Key_index[i+1]]
}
