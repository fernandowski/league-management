import StyledModal from "@/components/StyledModal";
import {Button, Text, TextInput} from "react-native-paper";
import {View, StyleSheet, useWindowDimensions, DimensionValue} from "react-native";
import {useEffect, useState} from "react";
import {MatchScore} from "@/hooks/useData";

interface ManageMatchModalProps {
    isOpen: boolean;
    onClose: () => void;
    onSave: (score: {home: number, away: number, matchId: string}) => void;
    matchDetails: MatchScore;
    seasonId: string;
}

export default function ManageMatchScoreModal(props: ManageMatchModalProps) {
    const [homeScore, setHomeScore] = useState(String(props.matchDetails.home_score));
    const [awayScore, setAwayScore] = useState(String(props.matchDetails.away_score));
    const dimensions = useWindowDimensions();
    const isLargeScreen = dimensions.width >= 768;

    const handleScoreChange = (
        value: string,
        setter: React.Dispatch<React.SetStateAction<string>>
    ) => {
        if (/^\d*$/.test(value)) {
            setter(value.replace(/^0+(?!$)/, ""));
        }
    };

    const modalWidth: DimensionValue = isLargeScreen ? "60%" : "90%";

    useEffect(() => {
        if (props.isOpen) {
            setHomeScore(props.matchDetails.home_score.toString());
            setAwayScore(props.matchDetails.away_score.toString());
        }
    }, [props.isOpen, props.matchDetails]);

    return (
        <StyledModal isOpen={props.isOpen} width={modalWidth}>
            <View style={{alignItems: "center"}}>
                <View style={[styles.header, !isLargeScreen ? {alignItems: "center"}: undefined]}>
                    <Text style={styles.title}>
                        {props.matchDetails.home_team}
                    </Text>
                    <Text style={styles.title}>
                        vs.
                    </Text>
                    <Text style={styles.title}>
                        {props.matchDetails.away_team}
                    </Text>
                </View>
                <View style={styles.body}>
                    <View style={[styles.scoreContainer, !isLargeScreen ? {flexDirection: "column"}: undefined]}>
                        <View style={styles.team}>
                            <Text style={styles.teamName}>{props.matchDetails.home_team}</Text>
                            <TextInput
                                style={styles.scoreInput}
                                value={homeScore}
                                keyboardType="numeric"
                                onChangeText={(value) => handleScoreChange(value, setHomeScore)}
                            />
                        </View>
                        <Text style={styles.separator}> - </Text>
                        <View style={styles.team}>
                            <Text style={styles.teamName}>{props.matchDetails.away_team}</Text>
                            <TextInput
                                style={styles.scoreInput}
                                value={awayScore}
                                keyboardType="numeric"
                                onChangeText={(value) => handleScoreChange(value, setAwayScore)}
                            />
                        </View>
                    </View>
                </View>
                <View style={styles.footer}>
                    <Button onPress={props.onClose}>Cancel</Button>
                    <Button
                        onPress={() => {
                            const parsedHomeScore = Number(homeScore);
                            const parsedAwayScore = Number(awayScore);
                            props.onSave({
                                matchId: props.matchDetails.id,
                                home: parsedHomeScore,
                                away: parsedAwayScore,
                            });
                        }}
                    >
                        Save
                    </Button>
                </View>
            </View>
        </StyledModal>
    );
}

const styles = StyleSheet.create({
    header: {
        marginBottom: 16,
    },
    title: {
        fontSize: 18,
        fontWeight: "bold",
        textAlign: "center",
    },
    body: {
        marginBottom: 16,
    },
    scoreContainer: {
        flexDirection: "row",
        alignItems: "center",
        justifyContent: "space-between",
    },
    team: {
        alignItems: "center",
    },
    teamName: {
        marginBottom: 8,
        width: 200,
        textAlign: "center"
    },
    scoreInput: {
        width: 60,
        height: 40,
        borderWidth: 1,
        borderColor: "#ddd",
        textAlign: "center",
        borderRadius: 4,
    },
    separator: {
        fontSize: 24,
        fontWeight: "bold",
    },
    footer: {
        flexDirection: "row",
        justifyContent: "space-between",
    },
});
