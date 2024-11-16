<?xml version="1.0" encoding="UTF-8" ?>
<!DOCTYPE mapper PUBLIC "-//mybatis.org//DTD Mapper 3.0//EN" "http://mybatis.org/dtd/mybatis-3-mapper.dtd" >
<mapper namespace="{{.DaoPackage}}">
    <!--auto generated Code-->
    <resultMap id="AllColumnMap" type="{{.PojoBean}}">
        {{.Results}}
    </resultMap>

    <sql id="all_column">
        {{.Columns}}
    </sql>

    <insert id="insertSelective">
        {{.InsertSQL}}
    </insert>

    <!--有条件插入-->
    <insert id="insert">
        INSERT INTO {{.InsertMap.table}}
        <trim prefix="(" suffix=")" suffixOverrides=",">
            {{.InsertMap.ifcolumn}}
        </trim>
        VALUES
        <trim prefix="(" suffix=")" suffixOverrides=",">
            {{.InsertMap.ifpojo}}
        </trim>
    </insert>

    <update id="updateOne">
        UPDATE {{.UpdateMap.table}}
        <trim prefix="set" suffixOverrides=",">
            {{.UpdateMap.updateSet}}
        </trim>
        WHERE sid=#{sid}
    </update>

    <select id="findList" parameterType="{{.PojoBean}}" resultMap="AllColumnMap">
        select
        <include refid="all_column"/>
        from {{.FindListMap.table}}
        <where>
            {{.FindListMap.where}}
        </where>
    </select>

    <delete id="delete" parameterType="{{.PojoBean}}">
        DELETE from {{.DeleteMap.table}}
        <where>
            {{.DeleteMap.where}}
        </where>
    </delete>


</mapper>