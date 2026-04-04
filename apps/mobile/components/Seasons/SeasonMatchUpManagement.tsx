import {StyleSheet, View} from "react-native";
import {Text} from "react-native-paper";
import {useCallback, useEffect, useState} from "react";
import {MatchesResponse, MatchScore, useData} from "@/hooks/useData";
import MatchUpRound from "@/components/Seasons/MatchUpRound";
import ManageMatchScoreModal from "@/components/Seasons/ManageMatchScoreModal";
import usePost from "@/hooks/usePost";

export interface SeasonMatchUpManagementProps {
    seasonId: string
}

export default function SeasonMatchUpManagement(props: SeasonMatchUpManagementProps) {
    const {fetchData, data, error} = useData<MatchesResponse[]>();
    const [openManageMatch, setOpenMangeMatch] = useState(false);
    const {postData} = usePost();
    const [matchDetails, setMatchDetails] = useState<MatchScore>();

    const onMatchPress = useCallback((matchScoreDetails: MatchScore) => {
        setMatchDetails(matchScoreDetails);
        setOpenMangeMatch(true);
    }, [])

    const onModalClose = () => {
        setOpenMangeMatch(false);
    }

    const onSave = async (scores: {home: number, away: number, matchId: string}) => {
        await postData(`/v1/seasons/${props.seasonId}/matches/score`, {
            match_id: scores.matchId,
            home_score: scores.home,
            away_score: scores.away,
        }, 'PUT');

        await fetchData(`/v1/seasons/${props.seasonId}/matches`);
        onModalClose();
    }

    useEffect(() => {
        if (props.seasonId) {
            fetchData(`/v1/seasons/${props.seasonId}/matches`)
        }

    }, [props.seasonId]);

    if (!data) {
        return <View><Text>{error}</Text></View>
    }

    return (
        <View>
            {
                data.sort((a, b) => (a.round - b.round)).map((match: MatchesResponse) => {
                    return (
                        <View style={[styles.matchUpRoundContainer]} key={match.round}>
                            <MatchUpRound data={match} onMatchPress={onMatchPress}/>
                        </View>
                    )
                })
            }
            { matchDetails && <ManageMatchScoreModal isOpen={openManageMatch} matchDetails={matchDetails} seasonId={props.seasonId} onSave={onSave} onClose={onModalClose}/>}
        </View>
    )
}

const styles = StyleSheet.create({
    matchUpRoundContainer: {
        marginTop: 18
    }
})
